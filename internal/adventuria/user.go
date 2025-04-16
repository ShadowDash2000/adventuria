package adventuria

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pocketbase/pocketbase/core"
)

type User struct {
	core.BaseRecordProxy
	LastAction Action
	Inventory  Inventory
	Timer      *Timer
	Stats      Stats
}

type Stats struct {
	Drops       int `json:"drops"`
	Rerolls     int `json:"rerolls"`
	Finished    int `json:"finished"`
	WasInJail   int `json:"wasInJail"`
	ItemsUsed   int `json:"itemsUsed"`
	DiceRolls   int `json:"diceRolls"`
	MaxDiceRoll int `json:"maxDiceRoll"`
	WheelRolled int `json:"wheelRolled"`
}

func NewUser(userId string) (*User, error) {
	if userId == "" {
		return nil, errors.New("you're not authorized")
	}

	var err error
	timer, err := NewTimer(userId)
	if err != nil {
		return nil, err
	}

	u := &User{
		Timer: timer,
	}

	err = u.fetchUser(userId)
	if err != nil {
		return nil, err
	}

	u.LastAction, err = NewLastUserAction(userId)
	if err != nil {
		return nil, err
	}

	u.Inventory, err = NewInventory(userId, u.MaxInventorySlots())
	if err != nil {
		return nil, err
	}

	u.bindHooks()

	return u, nil
}

func (u *User) bindHooks() {
	GameApp.OnRecordAfterUpdateSuccess(TableUsers).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == u.Id {
			u.SetProxyRecord(e.Record)
			u.Inventory.SetMaxSlots(u.MaxInventorySlots())
			u.UnmarshalJSONField("stats", &u.Stats)
		}
		return e.Next()
	})
}

func (u *User) fetchUser(userId string) error {
	user, err := GameApp.FindRecordById(TableUsers, userId)
	if err != nil {
		return err
	}

	u.SetProxyRecord(user)
	u.UnmarshalJSONField("stats", &u.Stats)

	return nil
}

func (u *User) CantDrop() bool {
	return u.GetBool("cantDrop")
}

func (u *User) SetCantDrop(b bool) {
	u.Set("cantDrop", b)
}

func (u *User) CanDrop() bool {
	return !u.CantDrop() && !u.IsInJail()
}

func (u *User) IsSafeDrop() bool {
	return u.DropsInARow() < GameSettings.DropsToJail()
}

func (u *User) IsInJail() bool {
	return u.GetBool("isInJail")
}

func (u *User) SetIsInJail(b bool) {
	u.Set("isInJail", b)
}

func (u *User) CurrentCell() (Cell, bool) {
	// TODO implement fake current cell

	cellsPassed := u.CellsPassed()
	currentCellNum := cellsPassed % GameCells.Count()

	return GameCells.GetByOrder(currentCellNum)
}

func (u *User) Points() int {
	return u.GetInt("points")
}

func (u *User) SetPoints(points int) {
	u.Set("points", points)
}

func (u *User) DropsInARow() int {
	return u.GetInt("dropsInARow")
}

func (u *User) SetDropsInARow(drops int) {
	u.Set("dropsInARow", drops)
}

func (u *User) CellsPassed() int {
	return u.GetInt("cellsPassed")
}

func (u *User) SetCellsPassed(cellsPassed int) {
	u.Set("cellsPassed", cellsPassed)
}

func (u *User) MaxInventorySlots() int {
	return u.GetInt("maxInventorySlots")
}

func (u *User) SetMaxInventorySlots(maxInventorySlots int) {
	u.Set("maxInventorySlots", maxInventorySlots)
}

func (u *User) ItemWheelsCount() int {
	return u.GetInt("itemWheelsCount")
}

func (u *User) SetItemWheelsCount(itemWheelsCount int) {
	u.Set("itemWheelsCount", itemWheelsCount)
}

func (u *User) Save() error {
	statsJson, _ := json.Marshal(u.Stats)
	u.Set("stats", string(statsJson))

	return GameApp.Save(u)
}

func (u *User) Move(steps int) (*OnAfterMoveFields, error) {
	cellsPassed := u.CellsPassed()
	currentCellNum := (cellsPassed + steps) % GameCells.Count()
	currentCell, _ := GameCells.GetByOrder(currentCellNum)

	u.SetCellsPassed(cellsPassed + steps)

	err := currentCell.OnCellReached(u)
	if err != nil {
		return nil, err
	}

	action := NewAction(u.Id, ActionTypeRollDice)
	action.SetCell(currentCell.ID())
	action.SetValue(steps)
	err = action.Save()
	if err != nil {
		return nil, err
	}

	prevCellNum := cellsPassed % GameCells.Count()
	lapsPassed := (prevCellNum + steps) / GameCells.Count()
	// Check if we're not moving backwards and passed new lap(-s)
	if steps > 0 && lapsPassed > 0 {
		onNewLapFields := &OnNewLapFields{
			Laps: lapsPassed,
		}
		GameEvent.Go(OnNewLap, NewEventFields(u, onNewLapFields))
	}

	onAfterMoveFields := &OnAfterMoveFields{
		Steps:       steps,
		Action:      action,
		CurrentCell: currentCell,
		Laps:        lapsPassed,
	}
	GameEvent.Go(OnAfterMove, NewEventFields(u, onAfterMoveFields))

	return onAfterMoveFields, nil
}

func (u *User) MoveToJail() error {
	jailCellPos, ok := GameCells.GetOrderByType(CellTypeJail)
	if !ok {
		return errors.New("jail cell not found")
	}

	currentCellNum := u.CellsPassed() % GameCells.Count()

	_, err := u.Move(jailCellPos - currentCellNum)
	if err != nil {
		return err
	}

	u.SetIsInJail(true)

	onAfterGoToJailFields := &OnAfterGoToJailFields{}
	GameEvent.Go(OnAfterGoToJail, NewEventFields(u, onAfterGoToJailFields))

	return nil
}

func (u *User) MoveToCellId(cellId string) error {
	cellPos, ok := GameCells.GetOrderById(cellId)
	if !ok {
		return fmt.Errorf("cell %s not found", cellId)
	}

	currentCellNum := u.CellsPassed() % GameCells.Count()

	_, err := u.Move(cellPos - currentCellNum)
	if err != nil {
		return err
	}

	return nil
}

// GetNextStepType
// WHAT IS THE NEXT STEP OF THE OPERATION? 👽
func (u *User) GetNextStepType() (string, error) {
	// Если еще не было сделано никаких lastAction, то делаем roll
	if u.LastAction == nil {
		return ActionTypeRollDice, nil
	}

	currentCell, ok := u.CurrentCell()
	if !ok {
		return "", errors.New("current cell not found")
	}

	lastActionType := ""
	if u.LastAction.CellId() == currentCell.ID() {
		lastActionType = u.LastAction.Type()
	}

	if currentCell.CantChooseAfterDrop() && lastActionType == ActionTypeDrop {
		return ActionTypeRollDice, nil
	}

	onBeforeNextStepFields := &OnBeforeNextStepFields{
		NextStepType: "",
		CurrentCell:  currentCell,
	}
	GameEvent.Go(OnBeforeNextStepType, NewEventFields(u, onBeforeNextStepFields))

	if onBeforeNextStepFields.NextStepType != "" {
		return onBeforeNextStepFields.NextStepType, nil
	}

	return currentCell.NextStep(u), nil
}

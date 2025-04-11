package adventuria

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pocketbase/pocketbase/core"
)

type User struct {
	core.BaseRecordProxy
	gc         *GameComponents
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

func NewUser(userId string, gc *GameComponents) (*User, error) {
	if userId == "" {
		return nil, errors.New("you're not authorized")
	}

	var err error
	timer, err := NewTimer(userId, gc)
	if err != nil {
		return nil, err
	}

	u := &User{
		gc:    gc,
		Timer: timer,
	}

	err = u.fetchUser(userId)
	if err != nil {
		return nil, err
	}

	u.LastAction, err = NewLastUserAction(userId, u.gc)
	if err != nil {
		return nil, err
	}

	u.Inventory, err = NewInventory(userId, u.MaxInventorySlots(), gc)
	if err != nil {
		return nil, err
	}

	u.bindHooks()

	return u, nil
}

func (u *User) bindHooks() {
	u.gc.App.OnRecordAfterUpdateSuccess(TableUsers).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == u.Id {
			u.SetProxyRecord(e.Record)
			u.Inventory.SetMaxSlots(u.MaxInventorySlots())
			u.UnmarshalJSONField("stats", &u.Stats)
		}
		return e.Next()
	})
}

func (u *User) fetchUser(userId string) error {
	user, err := u.gc.App.FindRecordById(TableUsers, userId)
	if err != nil {
		return err
	}

	u.SetProxyRecord(user)
	u.UnmarshalJSONField("stats", &u.Stats)

	return nil
}

func (u *User) IsSafeDrop() bool {
	return u.DropsInARow() < u.gc.Settings.DropsToJail()
}

func (u *User) SetIsSafeDrop(b bool) {
	u.Set("isSafeDrop", b)
}

func (u *User) IsInJail() bool {
	return u.GetBool("isInJail")
}

func (u *User) SetIsInJail(b bool) {
	u.Set("isInJail", b)
}

func (u *User) CurrentCell() (Cell, bool) {
	cellsPassed := u.CellsPassed()
	currentCellNum := cellsPassed % u.gc.Cells.Count()

	return u.gc.Cells.GetByOrder(currentCellNum)
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

	return u.gc.App.Save(u)
}

func (u *User) Move(steps int) (Action, Cell, error) {
	cellsPassed := u.CellsPassed()
	currentCellNum := (cellsPassed + steps) % u.gc.Cells.Count()
	currentCell, _ := u.gc.Cells.GetByOrder(currentCellNum)

	u.SetCellsPassed(cellsPassed + steps)

	err := currentCell.OnCellReached(u, u.gc)
	if err != nil {
		return nil, nil, err
	}

	action := NewAction(u.Id, ActionTypeRoll, u.gc)
	action.SetCell(currentCell.ID())
	action.SetValue(steps)
	err = action.Save()
	if err != nil {
		return nil, nil, err
	}

	prevCellNum := cellsPassed % u.gc.Cells.Count()
	lapsPassed := (prevCellNum + steps) / u.gc.Cells.Count()
	// Check if we're not moving backwards and passed new lap(-s)
	if steps > 0 && lapsPassed > 0 {
		onNewLapFields := &OnNewLapFields{
			Laps: lapsPassed,
		}
		u.gc.Event.Go(OnNewLap, NewEventFields(u, u.gc, onNewLapFields))
	}

	onAfterMoveFields := &OnAfterMoveFields{
		Steps:       steps,
		Action:      action,
		CurrentCell: currentCell,
	}
	u.gc.Event.Go(OnAfterMove, NewEventFields(u, u.gc, onAfterMoveFields))

	return action, currentCell, nil
}

func (u *User) MoveToJail() error {
	jailCellPos, ok := u.gc.Cells.GetOrderByType(CellTypeJail)
	if !ok {
		return errors.New("jail cell not found")
	}

	currentCellNum := u.CellsPassed() % u.gc.Cells.Count()

	_, _, err := u.Move(jailCellPos - currentCellNum)
	if err != nil {
		return err
	}

	u.SetIsInJail(true)

	onAfterGoToJailFields := &OnAfterGoToJailFields{}
	u.gc.Event.Go(OnAfterGoToJail, NewEventFields(u, u.gc, onAfterGoToJailFields))

	return nil
}

func (u *User) MoveToCellId(cellId string) error {
	cellPos, ok := u.gc.Cells.GetOrderById(cellId)
	if !ok {
		return fmt.Errorf("cell %s not found", cellId)
	}

	currentCellNum := u.CellsPassed() % u.gc.Cells.Count()

	_, _, err := u.Move(cellPos - currentCellNum)
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
		return ActionTypeRoll, nil
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
		return ActionTypeRoll, nil
	}

	onBeforeNextStepFields := &OnBeforeNextStepFields{
		NextStepType: "",
		CurrentCell:  currentCell,
	}
	u.gc.Event.Go(OnBeforeNextStepType, NewEventFields(u, u.gc, onBeforeNextStepFields))

	if onBeforeNextStepFields.NextStepType != "" {
		return onBeforeNextStepFields.NextStepType, nil
	}

	return currentCell.NextStep(u), nil
}

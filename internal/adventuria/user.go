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
	lastAction Action
	Inventory  *Inventory
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

	u.lastAction, err = NewLastUserAction(userId, u.gc)
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
	u.gc.app.OnRecordAfterUpdateSuccess(TableUsers).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == u.Id {
			u.SetProxyRecord(e.Record)
			u.Inventory.SetMaxSlots(u.MaxInventorySlots())
			u.UnmarshalJSONField("stats", &u.Stats)
		}
		return e.Next()
	})
}

func (u *User) fetchUser(userId string) error {
	user, err := u.gc.app.FindRecordById(TableUsers, userId)
	if err != nil {
		return err
	}

	u.SetProxyRecord(user)
	u.UnmarshalJSONField("stats", &u.Stats)

	return nil
}

func (u *User) IsSafeDrop() bool {
	return u.DropsInARow() < u.gc.settings.DropsToJail()
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

func (u *User) CurrentCell() (*Cell, bool) {
	cellsPassed := u.CellsPassed()
	currentCellNum := cellsPassed % u.gc.cells.Count()

	return u.gc.cells.GetByOrder(currentCellNum)
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

	return u.gc.app.Save(u)
}

func (u *User) Move(n int) (Action, *Cell, error) {
	cellsPassed := u.CellsPassed()
	currentCellNum := (cellsPassed + n) % u.gc.cells.Count()
	currentCell, _ := u.gc.cells.GetByOrder(currentCellNum)

	action := NewAction(u.Id, ActionTypeRoll, u.gc)
	action.SetCell(currentCell.Id)
	action.SetValue(n)
	err := action.Save()
	if err != nil {
		return nil, nil, err
	}

	u.SetCellsPassed(cellsPassed + n)

	prevCellNum := cellsPassed % u.gc.cells.Count()
	lapsPassed := (prevCellNum + n) / u.gc.cells.Count()
	// Check if we're not moving backwards and passed new lap(-s)
	if n > 0 && lapsPassed > 0 {
		u.gc.event.Go(OnNewLap, u, lapsPassed, u.gc)
	}

	return action, currentCell, nil
}

func (u *User) MoveToJail() error {
	jailCellPos, ok := u.gc.cells.GetOrderByType(CellTypeJail)
	if !ok {
		return errors.New("jail cell not found")
	}

	currentCellNum := u.CellsPassed() % u.gc.cells.Count()

	_, _, err := u.Move(jailCellPos - currentCellNum)
	if err != nil {
		return err
	}

	u.SetIsInJail(true)

	u.gc.event.Go(OnAfterGoToJail, u, u.gc)

	return nil
}

func (u *User) MoveToCellId(cellId string) error {
	cellPos, ok := u.gc.cells.GetOrderById(cellId)
	if !ok {
		return fmt.Errorf("cell %s not found", cellId)
	}

	currentCellNum := u.CellsPassed() % u.gc.cells.Count()

	_, _, err := u.Move(cellPos - currentCellNum)
	if err != nil {
		return err
	}

	return nil
}

// GetNextStepType
// WHAT IS THE NEXT STEP OF THE OPERATION? 👽
func (u *User) GetNextStepType() (string, error) {
	var nextStepType string

	// Если еще не было сделано никаких lastAction, то делаем roll
	if u.lastAction == nil {
		return ActionTypeRoll, nil
	}

	currentCell, ok := u.CurrentCell()
	if !ok {
		return nextStepType, errors.New("current cell not found")
	}

	lastActionType := ""
	if u.lastAction.CellId() == currentCell.Id {
		lastActionType = u.lastAction.Type()
	}

	if currentCell.CantChooseAfterDrop() && lastActionType == ActionTypeDrop {
		return ActionTypeRoll, nil
	}

	if u.ItemWheelsCount() > 0 {
		return ActionTypeRollItem, nil
	}

	// TODO: in future, maybe, this part needs to be in DB table
	switch currentCell.Type() {
	case CellTypeGame:
		switch lastActionType {
		case ActionTypeRoll,
			ActionTypeReroll:
			nextStepType = ActionTypeChooseGame
		case ActionTypeChooseGame:
			nextStepType = ActionTypeChooseResult
		case ActionTypeChooseResult,
			ActionTypeDrop:
			nextStepType = ActionTypeRoll
		default:
			nextStepType = ActionTypeChooseGame
		}
	case CellTypeStart:
		nextStepType = ActionTypeRoll
	case CellTypeJail:
		if u.IsInJail() {
			switch lastActionType {
			case ActionTypeRoll,
				ActionTypeReroll,
				ActionTypeDrop:
				nextStepType = ActionTypeRollCell
			case ActionTypeRollCell:
				nextStepType = ActionTypeChooseGame
			case ActionTypeChooseGame:
				nextStepType = ActionTypeChooseResult
			case ActionTypeChooseResult:
				nextStepType = ActionTypeRoll
			default:
				nextStepType = ActionTypeRollCell
			}
		} else {
			nextStepType = ActionTypeRoll
		}
	case CellTypePreset:
		switch lastActionType {
		case ActionTypeRoll:
			nextStepType = ActionTypeRollWheelPreset
		case ActionTypeReroll,
			ActionTypeDrop:
			nextStepType = ActionTypeRollWheelPreset
		case ActionTypeRollWheelPreset:
			nextStepType = ActionTypeChooseGame
		case ActionTypeChooseGame:
			nextStepType = ActionTypeChooseResult
		case ActionTypeChooseResult:
			nextStepType = ActionTypeRoll
		default:
			nextStepType = ActionTypeRollWheelPreset
		}
	case CellTypeItem:
		switch lastActionType {
		case ActionTypeRoll:
			nextStepType = ActionTypeRollItem
		case ActionTypeRollItem:
			nextStepType = ActionTypeRoll
		default:
			nextStepType = ActionTypeRollItem
		}
	case CellTypeWheelPreset:
		switch lastActionType {
		case ActionTypeRoll,
			ActionTypeReroll:
			nextStepType = ActionTypeRollWheelPreset
		case ActionTypeRollWheelPreset:
			nextStepType = ActionTypeChooseResult
		case ActionTypeChooseResult,
			ActionTypeDrop:
			nextStepType = ActionTypeRoll
		default:
			nextStepType = ActionTypeRollWheelPreset
		}
	}

	return nextStepType, nil
}

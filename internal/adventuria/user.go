package adventuria

import (
	"adventuria/pkg/collections"
	"encoding/json"
	"errors"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type User struct {
	app        core.App
	log        *Log
	userId     string
	user       *core.Record
	lastAction *core.Record
	Inventory  *Inventory
	cells      *Cells
	settings   *Settings
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

func NewUser(
	userId string,
	cells *Cells,
	settings *Settings,
	log *Log, cols *collections.Collections,
	app core.App,
) (*User, error) {
	if userId == "" {
		return nil, errors.New("you're not authorized")
	}

	var err error
	timer, err := NewTimer(userId, settings, cols, app)
	if err != nil {
		return nil, err
	}

	u := &User{
		app:      app,
		log:      log,
		userId:   userId,
		cells:    cells,
		settings: settings,
		Timer:    timer,
	}

	err = u.fetchUser()
	if err != nil {
		return nil, err
	}

	err = u.fetchUserAction()
	if err != nil {
		return nil, err
	}

	u.Inventory, err = NewInventory(userId, u.user.GetInt("maxInventorySlots"), log, cols, app)
	if err != nil {
		return nil, err
	}

	u.bindHooks()

	return u, nil
}

func (u *User) bindHooks() {
	u.app.OnRecordAfterCreateSuccess(TableActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")
		if userId == u.userId {
			u.lastAction = e.Record
		}
		return e.Next()
	})
	u.app.OnRecordAfterUpdateSuccess(TableActions).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == u.lastAction.Id {
			u.lastAction = e.Record
		}
		return e.Next()
	})
	u.app.OnRecordAfterDeleteSuccess(TableActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")
		if userId == u.userId {
			u.fetchUserAction()
		}
		return e.Next()
	})
	u.app.OnRecordAfterUpdateSuccess(TableUsers).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == u.userId {
			u.user = e.Record
			u.Inventory.SetMaxSlots(e.Record.GetInt("maxInventorySlots"))
			u.user.UnmarshalJSONField("stats", &u.Stats)
		}
		return e.Next()
	})
}

func (u *User) fetchUser() error {
	var err error
	u.user, err = u.app.FindRecordById(TableUsers, u.userId)
	if err != nil {
		return err
	}

	u.user.UnmarshalJSONField("stats", &u.Stats)

	return nil
}

func (u *User) fetchUserAction() error {
	actions, err := u.app.FindRecordsByFilter(
		TableActions,
		"user.id = {:userId}",
		"-created",
		1,
		0,
		dbx.Params{"userId": u.userId},
	)
	if err != nil {
		return err
	}

	if len(actions) > 0 {
		u.lastAction = actions[0]
	}

	return nil
}

func (u *User) IsSafeDrop() bool {
	return u.DropsInARow() < u.settings.DropsToJail()
}

func (u *User) IsInJail() bool {
	return u.user.GetBool("isInJail")
}

func (u *User) CurrentCell() (*core.Record, bool) {
	cellsPassed := u.CellsPassed()
	currentCellNum := cellsPassed % u.cells.Count()

	return u.cells.GetByOrder(currentCellNum)
}

func (u *User) Points() int {
	return u.user.GetInt("points")
}

func (u *User) DropsInARow() int {
	return u.user.GetInt("dropsInARow")
}

func (u *User) CellsPassed() int {
	return u.user.GetInt("cellsPassed")
}

func (u *User) ItemWheelsCount() int {
	return u.user.GetInt("itemWheelsCount")
}

func (u *User) Set(key string, value any) {
	u.user.Set(key, value)
}

func (u *User) Save() error {
	statsJson, _ := json.Marshal(u.Stats)
	u.user.Set("stats", string(statsJson))

	return u.app.Save(u.user)
}

// GetNextStepType
// WHAT IS THE NEXT STEP OF THE OPERATION? 👽
func (u *User) GetNextStepType() (string, error) {
	var nextStepType string

	// Если еще не было сделано никаких lastAction, то делаем roll
	if u.lastAction == nil {
		return ActionTypeRoll, nil
	}

	cell, ok := u.CurrentCell()
	if !ok {
		return nextStepType, errors.New("current cell not found")
	}

	cellType := cell.GetString("type")
	cantChooseAfterDrop := cell.GetBool("cantChooseAfterDrop")
	lastActionType := ""
	if u.lastAction.GetString("cell") == cell.Id {
		lastActionType = u.lastAction.GetString("type")
	}

	if cantChooseAfterDrop && lastActionType == ActionTypeDrop {
		return ActionTypeRoll, nil
	}

	if u.ItemWheelsCount() > 0 {
		return ActionTypeRollItem, nil
	}

	// TODO: in future, maybe, this part needs to be in DB table
	switch cellType {
	case CellTypeGame:
		switch lastActionType {
		case ActionTypeRoll,
			ActionTypeReroll,
			ActionTypeDrop:
			nextStepType = ActionTypeChooseGame
		case ActionTypeChooseGame:
			nextStepType = ActionTypeChooseResult
		case ActionTypeChooseResult:
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
			ActionTypeReroll,
			ActionTypeDrop:
			nextStepType = ActionTypeRollWheelPreset
		case ActionTypeRollWheelPreset:
			nextStepType = ActionTypeChooseResult
		case ActionTypeChooseResult:
			nextStepType = ActionTypeRoll
		default:
			nextStepType = ActionTypeRollWheelPreset
		}
	}

	return nextStepType, nil
}

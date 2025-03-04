package adventuria

import (
	"adventuria/pkg/collections"
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
		app:    app,
		log:    log,
		userId: userId,
		cells:  cells,
		Timer:  timer,
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
	return u.GetDropsInARow() < 2
}

func (u *User) IsInJail() bool {
	return u.user.GetBool("isInJail")
}

func (u *User) GetCurrentCell() (*core.Record, bool) {
	cellsPassed := u.GetCellsPassed()
	currentCellNum := cellsPassed % u.cells.Count()

	return u.cells.GetByOrder(currentCellNum)
}

func (u *User) GetPoints() int {
	return u.user.GetInt("points")
}

func (u *User) GetDropsInARow() int {
	return u.user.GetInt("dropsInARow")
}

func (u *User) GetCellsPassed() int {
	return u.user.GetInt("cellsPassed")
}

func (u *User) Set(key string, value any) {
	u.user.Set(key, value)
}

func (u *User) Save() error {
	return u.app.Save(u.user)
}

// GetNextStepType
// WHAT IS THE NEXT STEP OF THE OPERATION? 👽
func (u *User) GetNextStepType() (string, error) {
	var nextStepType string

	// Если еще не было сделано никаких lastAction, то делаем roll
	if u.lastAction == nil {
		return UserNextStepRoll, nil
	}

	cell, ok := u.GetCurrentCell()
	if !ok {
		return nextStepType, errors.New("current cell not found")
	}

	cellType := cell.GetString("type")
	lastActionType := ""
	if u.lastAction.GetString("cell") == cell.Id {
		lastActionType = u.lastAction.GetString("type")
	}

	switch cellType {
	case CellTypeGame:
		switch lastActionType {
		case ActionTypeRoll,
			ActionTypeReroll,
			ActionTypeDrop:
			nextStepType = UserNextStepChooseGame
		case ActionTypeGame:
			nextStepType = UserNextStepChooseResult
		case ActionTypeDone:
			nextStepType = UserNextStepRoll
		default:
			nextStepType = UserNextStepChooseGame
		}
	case CellTypeStart:
		nextStepType = UserNextStepRoll
	case CellTypeJail:
		if u.IsInJail() {
			switch lastActionType {
			case ActionTypeRoll:
				nextStepType = UserNextStepRollJailCell
			case ActionTypeReroll,
				ActionTypeDrop,
				ActionTypeRollCell:
				nextStepType = UserNextStepChooseGame
			case ActionTypeGame:
				nextStepType = UserNextStepChooseResult
			case ActionTypeDone:
				nextStepType = UserNextStepRoll
			default:
				nextStepType = UserNextStepRoll
			}
		} else {
			nextStepType = UserNextStepRoll
		}
	case CellTypeBigWin:
		switch lastActionType {
		case ActionTypeRoll,
			ActionTypeReroll:
			nextStepType = UserNextStepRollBigWin
		case ActionTypeDrop:
			nextStepType = UserNextStepRoll
		case ActionTypeGame,
			ActionTypeRollBigWin:
			nextStepType = UserNextStepChooseResult
		case ActionTypeDone:
			nextStepType = UserNextStepRoll
		default:
			nextStepType = UserNextStepRollBigWin
		}
	case CellTypePreset:
		switch lastActionType {
		case ActionTypeRoll:
			nextStepType = UserNextStepRollPreset
		case ActionTypeReroll,
			ActionTypeDrop,
			ActionTypeRollPreset:
			nextStepType = UserNextStepChooseGame
		case ActionTypeGame:
			nextStepType = UserNextStepChooseResult
		case ActionTypeDone:
			nextStepType = UserNextStepRoll
		default:
			nextStepType = UserNextStepRollPreset
		}
	case CellTypeMovie:
		switch lastActionType {
		case ActionTypeRoll:
			nextStepType = UserNextStepRollMovie
		case ActionTypeRollMovie:
			nextStepType = UserNextStepMovieResult
		case ActionTypeDone:
		case ActionTypeMovieResult:
			nextStepType = UserNextStepRoll
		default:
			nextStepType = UserNextStepRollMovie
		}
	case CellTypeItem:
		switch lastActionType {
		case ActionTypeRoll:
			nextStepType = UserNextStepRollItem
		case ActionTypeRollItem:
			nextStepType = UserNextStepRoll
		default:
			nextStepType = UserNextStepRollItem
		}
	case CellTypeDeveloper:
		switch lastActionType {
		case ActionTypeRoll,
			ActionTypeReroll,
			ActionTypeDrop:
			nextStepType = UserNextStepRollDeveloper
		case ActionTypeRollDeveloper:
			nextStepType = UserNextStepChooseResult
		case ActionTypeDone:
			nextStepType = UserNextStepRoll
		default:
			nextStepType = UserNextStepRollDeveloper
		}
	}

	return nextStepType, nil
}

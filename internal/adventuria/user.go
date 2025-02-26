package adventuria

import (
	"errors"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type User struct {
	app        core.App
	userId     string
	user       *core.Record
	lastAction *core.Record
	Inventory  *Inventory
}

func NewUser(userId string, app core.App) (*User, error) {
	if userId == "" {
		return nil, errors.New("you're not authorized")
	}

	u := &User{
		app:    app,
		userId: userId,
	}

	var err error

	err = u.fetchUser()
	if err != nil {
		return nil, err
	}

	err = u.fetchUserAction()
	if err != nil {
		return nil, err
	}

	u.Inventory, err = NewInventory(userId, u.user.GetInt("maxInventorySlots"), app)
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

func (u *User) GetCurrentCell() (*core.Record, error) {
	if u.lastAction == nil {
		return nil, errors.New("no lastAction found")
	}

	errs := u.app.ExpandRecord(u.lastAction, []string{"cell"}, nil)
	if len(errs) > 0 {
		for _, err := range errs {
			return nil, err
		}
	}

	return u.lastAction.ExpandedOne("cell"), nil
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

	cell, err := u.GetCurrentCell()
	if err != nil {
		return nextStepType, err
	}

	cellType := cell.GetString("type")
	lastActionType := u.lastAction.GetString("type")

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
			nextStepType = UserNextStepRollOnBigWinDrop
		case ActionTypeGame,
			ActionTypeRollBigWin:
			nextStepType = UserNextStepChooseResult
		case ActionTypeDone:
			nextStepType = UserNextStepRoll
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
		}
	case CellTypeMovie:
		switch lastActionType {
		case ActionTypeRoll:
			nextStepType = UserNextStepRollMovie
		case ActionTypeRollMovie:
			nextStepType = UserNextStepMovieResult
		}
	case CellTypeItem:
		switch lastActionType {
		case ActionTypeRoll:
			nextStepType = UserNextStepRollItem
		case ActionTypeRollItem:
			nextStepType = UserNextStepRoll
		}
	}

	return nextStepType, nil
}

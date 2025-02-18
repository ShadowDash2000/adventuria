package usecases

import (
	"adventuria/internal/adventuria"
	"errors"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"slices"
)

type User struct {
	app       core.App
	userId    string
	user      *core.Record
	actions   []*core.Record
	Inventory *Inventory
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

	err = u.fetchUserActions(3)
	if err != nil {
		return nil, err
	}

	u.Inventory, err = NewInventory(userId, app)
	if err != nil {
		return nil, err
	}

	u.bindHooks()

	return u, nil
}

func (u *User) bindHooks() {
	u.app.OnRecordAfterCreateSuccess(adventuria.TableActions).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString("user")
		if userId == u.userId {
			u.addAction(e.Record)
		}
		return e.Next()
	})

	u.app.OnRecordAfterUpdateSuccess(adventuria.TableUsers).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == u.userId {
			u.user = e.Record
		}
		return e.Next()
	})
}

func (u *User) fetchUser() error {
	var err error
	u.user, err = u.app.FindRecordById(adventuria.TableUsers, u.userId)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) fetchUserActions(limit int) error {
	var err error
	u.actions, err = u.app.FindRecordsByFilter(
		adventuria.TableActions,
		"user.id = {:userId}",
		"-created",
		limit,
		0,
		dbx.Params{"userId": u.userId},
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) addAction(action *core.Record) {
	copy(u.actions[1:], u.actions)
	u.actions[0] = action
}

func (u *User) CanRoll() (bool, error) {
	cell, err := u.GetCurrentCell()
	if err != nil {
		return false, err
	}

	action := u.actions[0]
	actionFields := action.FieldsData()
	cellFields := cell.FieldsData()
	canRoll := true

	statuses := []string{
		adventuria.ActionStatusGameNotChosen,
		adventuria.ActionStatusReroll,
		adventuria.ActionStatusDrop,
		adventuria.ActionStatusInProgress,
	}
	if slices.Contains(statuses, actionFields["status"].(string)) {
		canRoll = false
	}

	if cellFields["code"].(string) == adventuria.CellTypeBigWin &&
		actionFields["status"].(string) == adventuria.ActionStatusReroll {
		canRoll = true
	}

	return canRoll, nil
}

func (u *User) CanDrop() (bool, error) {
	if len(u.actions) == 0 {
		return false, nil
	}

	if len(u.actions) == 3 {
		previousActions := u.actions[1:3]
		i := 0

		for _, previousAction := range previousActions {
			previousActionFields := previousAction.FieldsData()

			if previousActionFields["status"].(string) == adventuria.ActionStatusDrop {
				i++
			}
		}

		if i >= 2 {
			return false, nil
		}
	}

	action := u.actions[0]
	actionFields := action.FieldsData()

	if actionFields["status"].(string) != adventuria.ActionStatusInProgress {
		return false, nil
	}

	return true, nil
}

func (u *User) IsInJail() (bool, error) {
	if len(u.actions) == 0 {
		return false, nil
	}

	cell, err := u.GetCurrentCell()
	if err != nil {
		return false, err
	}
	cellFields := cell.FieldsData()

	if cellFields["type"].(string) != adventuria.CellTypeJail {
		return false, nil
	}

	if len(u.actions) == 3 {
		previousActions := u.actions[1:3]
		i := 0

		for _, previousAction := range previousActions {
			previousActionFields := previousAction.FieldsData()

			if previousActionFields["status"].(string) == adventuria.ActionStatusDrop {
				i++
			}
		}

		if i >= 2 {
			return true, nil
		}
	}

	return false, nil
}

func (u *User) GetCurrentCell() (*core.Record, error) {
	if len(u.actions) == 0 {
		return nil, errors.New("no actions found")
	}

	action := u.actions[0]

	errs := u.app.ExpandRecord(action, []string{"cell"}, nil)
	if len(errs) > 0 {
		for _, err := range errs {
			return nil, err
		}
	}

	return action.ExpandedOne("cell"), nil
}

func (u *User) GetPoints() int {
	return u.user.GetInt("points")
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

	// Если еще не было сделано никаких actions, то делаем roll
	if len(u.actions) == 0 {
		return adventuria.UserNextStepRoll, nil
	}

	cell, err := u.GetCurrentCell()
	if err != nil {
		return nextStepType, err
	}
	cellFields := cell.FieldsData()
	action := u.actions[0]
	actionFields := action.FieldsData()

	cellType := cellFields["type"].(string)
	actionStatus := actionFields["status"].(string)

	switch cellType {
	case adventuria.CellTypeGame:
		switch actionStatus {
		case adventuria.ActionStatusNone,
			adventuria.ActionStatusDone:
			nextStepType = adventuria.UserNextStepRoll
		case adventuria.ActionStatusGameNotChosen,
			adventuria.ActionStatusReroll,
			adventuria.ActionStatusDrop:
			nextStepType = adventuria.UserNextStepChooseGame
		case adventuria.ActionStatusInProgress:
			nextStepType = adventuria.UserNextStepChooseResult
		}
	case adventuria.CellTypeStart:
		switch actionStatus {
		case adventuria.ActionStatusNone,
			adventuria.ActionStatusDone,
			adventuria.ActionStatusGameNotChosen,
			adventuria.ActionStatusReroll,
			adventuria.ActionStatusDrop,
			adventuria.ActionStatusInProgress:
			nextStepType = adventuria.UserNextStepRoll
		}
	case adventuria.CellTypeJail:
		isInJail, err := u.IsInJail()
		if err != nil {
			return nextStepType, err
		}

		if isInJail {
			switch actionStatus {
			case adventuria.ActionStatusNone:
			case adventuria.ActionStatusDone:
				nextStepType = adventuria.UserNextStepRoll
			case adventuria.ActionStatusGameNotChosen:
				nextStepType = adventuria.UserNextStepRollCell
			case adventuria.ActionStatusReroll,
				adventuria.ActionStatusDrop:
				nextStepType = adventuria.UserNextStepChooseGame
			case adventuria.ActionStatusInProgress:
				nextStepType = adventuria.UserNextStepChooseResult
			}
		} else {
			switch actionStatus {
			case adventuria.ActionStatusNone,
				adventuria.ActionStatusDone,
				adventuria.ActionStatusGameNotChosen,
				adventuria.ActionStatusReroll,
				adventuria.ActionStatusInProgress:
				nextStepType = adventuria.UserNextStepRoll
			}
		}
	case adventuria.CellTypeBigWin:
		switch actionStatus {
		case adventuria.ActionStatusNone,
			adventuria.ActionStatusDone,
			adventuria.ActionStatusGameNotChosen,
			adventuria.ActionStatusReroll,
			adventuria.ActionStatusInProgress:
			nextStepType = adventuria.UserNextStepRoll
		}
	case adventuria.CellTypePreset:
		switch actionStatus {
		case adventuria.ActionStatusNone,
			adventuria.ActionStatusDone,
			adventuria.ActionStatusGameNotChosen,
			adventuria.ActionStatusReroll,
			adventuria.ActionStatusInProgress:
			nextStepType = adventuria.UserNextStepRoll
		}
	}

	return nextStepType, nil
}

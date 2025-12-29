package adventuria

import (
	"errors"
)

type ActionBase struct {
	t    ActionType
	user User
}

func NewActionFromType(actionType ActionType) (Action, error) {
	actionCreator, ok := actionsList[actionType]
	if !ok {
		return nil, errors.New("unknown action type")
	}

	action := actionCreator()

	return action, nil
}

func (a *ActionBase) Type() ActionType {
	return a.t
}

func (a *ActionBase) setType(t ActionType) {
	a.t = t
}

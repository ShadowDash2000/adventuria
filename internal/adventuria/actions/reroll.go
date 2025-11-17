package actions

import (
	"adventuria/internal/adventuria"
	"errors"
)

type RerollAction struct {
	adventuria.ActionBase
}

func (a *RerollAction) CanDo(user adventuria.User) bool {
	currentCell, ok := user.CurrentCell()
	if ok {
		if currentCell.CantReroll() {
			return false
		}
	}

	return user.LastAction().Type() == ActionTypeRollWheel
}

func (a *RerollAction) Do(user adventuria.User, req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	var comment string
	if c, ok := req["comment"]; ok {
		comment = c.(string)
	}

	currentCell, ok := user.CurrentCell()
	if !ok {
		return nil, errors.New("current cell not found")
	}

	if currentCell.CantReroll() {
		return nil, errors.New("can't reroll on this cell")
	}

	action := user.LastAction()
	action.SetType(ActionTypeReroll)
	action.SetComment(comment)

	err := user.OnAfterReroll().Trigger(&adventuria.OnAfterRerollEvent{})
	if err != nil {
		return nil, err
	}

	return &adventuria.ActionResult{}, nil
}

package actions

import (
	"adventuria/internal/adventuria"
	"errors"
)

type RerollAction struct {
	adventuria.ActionBase
}

func (a *RerollAction) CanDo() bool {
	currentCell, ok := a.User().CurrentCell()
	if ok {
		if currentCell.CantReroll() {
			return false
		}
	}

	return a.User().LastAction().Type() == ActionTypeRollWheel
}

func (a *RerollAction) NextAction() adventuria.ActionType {
	return ActionTypeRollWheel
}

func (a *RerollAction) Do(req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	currentCell, ok := a.User().CurrentCell()
	if !ok {
		return nil, errors.New("current cell not found")
	}

	if currentCell.CantReroll() {
		return nil, errors.New("can't reroll on this cell")
	}

	action := a.User().LastAction()
	action.SetType(ActionTypeReroll)
	action.SetComment(req.Comment)

	err := a.User().OnAfterReroll().Trigger(&adventuria.OnAfterRerollEvent{})
	if err != nil {
		return nil, err
	}

	return &adventuria.ActionResult{}, nil
}

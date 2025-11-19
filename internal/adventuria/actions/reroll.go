package actions

import (
	"adventuria/internal/adventuria"
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
		comment, ok = c.(string)
		if !ok {
			return &adventuria.ActionResult{
				Success: false,
				Error:   "request error: comment is not string",
			}, nil
		}
	}

	action := user.LastAction()
	action.SetType(ActionTypeReroll)
	action.SetComment(comment)

	err := user.OnAfterReroll().Trigger(&adventuria.OnAfterRerollEvent{})
	if err != nil {
		adventuria.PocketBase.Logger().Error(
			"reroll.do(): failed to trigger onAfterReroll event",
			"error",
			err,
		)
	}

	return &adventuria.ActionResult{
		Success: true,
	}, nil
}

package actions

import (
	"adventuria/internal/adventuria"
	"fmt"
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
	err := adventuria.PocketBase.Save(action.ProxyRecord())
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't save action record",
		}, fmt.Errorf("reroll.do(): %w", err)
	}
	action.ProxyRecord().MarkAsNew()
	action.ProxyRecord().Set("id", "")
	action.SetComment("")
	action.SetGame("")
	action.SetDiceRoll(0)

	err = user.OnAfterReroll().Trigger(&adventuria.OnAfterRerollEvent{})
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

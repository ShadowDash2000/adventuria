package actions

import (
	"adventuria/internal/adventuria"
	"errors"
)

type DoneAction struct {
	adventuria.ActionBase
}

func (a *DoneAction) CanDo(user adventuria.User) bool {
	return user.LastAction().Type() == ActionTypeRollWheel
}

func (a *DoneAction) Do(user adventuria.User, req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
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

	currentCell, ok := user.CurrentCell()
	if !ok {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: current cell not found",
		}, errors.New("done.do(): current cell not found")
	}

	onBeforeDoneEvent := &adventuria.OnBeforeDoneEvent{
		CellPoints: currentCell.Points(),
	}
	err := user.OnBeforeDone().Trigger(onBeforeDoneEvent)
	if err != nil {
		adventuria.PocketBase.Logger().Error(
			"done.do(): failed to trigger onBeforeDone event",
			"error",
			err,
		)
	}

	action := user.LastAction()
	action.SetType(ActionTypeDone)
	action.SetComment(comment)
	action.SetCanMove(true)

	user.SetDropsInARow(0)
	user.SetIsInJail(false)
	user.SetPoints(user.Points() + onBeforeDoneEvent.CellPoints)

	onAfterDoneEvent := &adventuria.OnAfterDoneEvent{}
	err = user.OnAfterDone().Trigger(onAfterDoneEvent)
	if err != nil {
		adventuria.PocketBase.Logger().Error(
			"done.do(): failed to trigger onAfterDone event",
			"error",
			err,
		)
	}

	return &adventuria.ActionResult{
		Success: true,
	}, nil
}

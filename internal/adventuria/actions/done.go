package actions

import (
	"adventuria/internal/adventuria"
	"errors"
	"fmt"
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

	// TODO: we need to pass actual cell points here instead of doing math here
	onBeforeDoneEvent := &adventuria.OnBeforeDoneEvent{
		CellPointsDivide: 0,
	}
	err := user.OnBeforeDone().Trigger(onBeforeDoneEvent)
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error",
		}, fmt.Errorf("done.do(): %w", err)
	}

	action := user.LastAction()
	action.SetType(ActionTypeDone)
	action.SetComment(comment)
	action.SetCanMove(true)

	cellPoints := currentCell.Points()
	if onBeforeDoneEvent.CellPointsDivide != 0 {
		cellPoints /= onBeforeDoneEvent.CellPointsDivide
	}

	user.SetDropsInARow(0)
	user.SetIsInJail(false)
	user.SetPoints(user.Points() + cellPoints)

	onAfterDoneEvent := &adventuria.OnAfterDoneEvent{}
	err = user.OnAfterDone().Trigger(onAfterDoneEvent)
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error",
		}, fmt.Errorf("done.do(): %w", err)
	}

	return &adventuria.ActionResult{
		Success: true,
	}, nil
}

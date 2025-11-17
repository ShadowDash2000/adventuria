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
		comment = c.(string)
	}

	currentCell, ok := user.CurrentCell()
	if !ok {
		return nil, errors.New("current cell not found")
	}

	onBeforeDoneEvent := &adventuria.OnBeforeDoneEvent{
		CellPointsDivide: 0,
	}
	err := user.OnBeforeDone().Trigger(onBeforeDoneEvent)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &adventuria.ActionResult{
		Success: true,
	}, nil
}

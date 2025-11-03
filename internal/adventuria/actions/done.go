package actions

import (
	"adventuria/internal/adventuria"
	"errors"
)

type DoneAction struct {
	adventuria.ActionBase
}

func (a *DoneAction) CanDo() bool {
	return a.User().LastAction().Type() == ActionTypeRollWheel
}

func (a *DoneAction) Do(req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	var comment string
	if c, ok := req["comment"]; ok {
		comment = c.(string)
	}

	currentCell, ok := a.User().CurrentCell()
	if !ok {
		return nil, errors.New("current cell not found")
	}

	onBeforeDoneEvent := &adventuria.OnBeforeDoneEvent{
		CellPointsDivide: 0,
	}
	err := a.User().OnBeforeDone().Trigger(onBeforeDoneEvent)
	if err != nil {
		return nil, err
	}

	action := a.User().LastAction()
	action.SetType(ActionTypeDone)
	action.SetComment(comment)
	action.SetCanMove(true)

	cellPoints := currentCell.Points()
	if onBeforeDoneEvent.CellPointsDivide != 0 {
		cellPoints /= onBeforeDoneEvent.CellPointsDivide
	}

	a.User().SetDropsInARow(0)
	a.User().SetIsInJail(false)
	a.User().SetPoints(a.User().Points() + cellPoints)

	onAfterDoneEvent := &adventuria.OnAfterDoneEvent{}
	err = a.User().OnAfterDone().Trigger(onAfterDoneEvent)
	if err != nil {
		return nil, err
	}

	return &adventuria.ActionResult{
		Success: true,
	}, nil
}

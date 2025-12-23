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
		CellCoins:  currentCell.Coins(),
	}
	res, err := user.OnBeforeDone().Trigger(onBeforeDoneEvent)
	if res != nil && !res.Success {
		return &adventuria.ActionResult{
			Success: false,
			Error:   res.Error,
		}, err
	}
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: failed to trigger onBeforeDone event",
		}, err
	}

	action := user.LastAction()
	action.SetType(ActionTypeDone)
	action.SetComment(comment)
	action.SetCanMove(true)

	user.SetDropsInARow(0)
	user.SetIsInJail(false)
	user.SetPoints(user.Points() + onBeforeDoneEvent.CellPoints)
	user.SetBalance(user.Balance() + onBeforeDoneEvent.CellCoins)

	onAfterDoneEvent := &adventuria.OnAfterDoneEvent{}
	res, err = user.OnAfterDone().Trigger(onAfterDoneEvent)
	if res != nil && !res.Success {
		return &adventuria.ActionResult{
			Success: false,
			Error:   res.Error,
		}, err
	}
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: failed to trigger onAfterDoneEvent event",
		}, err
	}

	return &adventuria.ActionResult{
		Success: true,
	}, nil
}

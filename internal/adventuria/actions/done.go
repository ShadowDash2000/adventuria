package actions

import (
	"adventuria/internal/adventuria"
	"errors"
	"fmt"
)

type DoneAction struct {
	adventuria.ActionBase
}

func (a *DoneAction) CanDo(ctx adventuria.ActionContext) bool {
	return ctx.User.LastAction().Type() == ActionTypeRollWheel
}

func (a *DoneAction) Do(ctx adventuria.ActionContext, req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
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

	currentCell, ok := ctx.User.CurrentCell()
	if !ok {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: current cell not found",
		}, errors.New("done.do(): current cell not found")
	}

	onBeforeDoneEvent := &adventuria.OnBeforeDoneEvent{
		AppContext: ctx.AppContext,
		CellPoints: currentCell.Points(),
		CellCoins:  currentCell.Coins(),
	}
	res, err := ctx.User.OnBeforeDone().Trigger(onBeforeDoneEvent)
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

	action := ctx.User.LastAction()
	action.SetType(ActionTypeDone)
	action.SetComment(comment)
	action.SetCanMove(true)

	ctx.User.SetDropsInARow(0)
	ctx.User.SetIsInJail(false)
	ctx.User.SetPoints(ctx.User.Points() + onBeforeDoneEvent.CellPoints)
	err = ctx.User.AddBalance(ctx.AppContext, onBeforeDoneEvent.CellCoins)
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't update user balance",
		}, fmt.Errorf("done.do(): %w", err)
	}

	res, err = ctx.User.OnAfterDone().Trigger(&adventuria.OnAfterDoneEvent{
		AppContext: ctx.AppContext,
	})
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

func (a *DoneAction) GetVariants(_ adventuria.ActionContext) any {
	return nil
}

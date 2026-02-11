package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/result"
	"errors"
)

type DoneAction struct {
	adventuria.ActionBase
}

func (a *DoneAction) CanDo(ctx adventuria.ActionContext) bool {
	return ctx.User.LastAction().Type() == ActionTypeRollWheel
}

func (a *DoneAction) Do(ctx adventuria.ActionContext, req adventuria.ActionRequest) (*result.Result, error) {
	var comment string
	if c, ok := req["comment"]; ok {
		comment, ok = c.(string)
		if !ok {
			return result.Err("comment is not string"), nil
		}
	}

	currentCell, ok := ctx.User.CurrentCell()
	if !ok {
		return result.Err("internal error: current cell not found"),
			errors.New("done.do(): current cell not found")
	}

	onBeforeDoneEvent := &adventuria.OnBeforeDoneEvent{
		AppContext: ctx.AppContext,
		CellPoints: currentCell.Points(),
		CellCoins:  currentCell.Coins(),
	}
	res, err := ctx.User.OnBeforeDone().Trigger(onBeforeDoneEvent)
	if err != nil {
		return res, err
	}
	if res.Failed() {
		return res, err
	}

	action := ctx.User.LastAction()
	action.SetType(ActionTypeDone)
	action.SetComment(comment)
	action.SetCanMove(true)

	ctx.User.SetDropsInARow(0)
	ctx.User.SetIsInJail(false)
	ctx.User.SetPoints(ctx.User.Points() + onBeforeDoneEvent.CellPoints)
	ctx.User.AddBalance(onBeforeDoneEvent.CellCoins)

	res, err = ctx.User.OnAfterDone().Trigger(&adventuria.OnAfterDoneEvent{
		AppContext: ctx.AppContext,
	})
	if err != nil {
		return res, err
	}
	if res.Failed() {
		return res, err
	}

	return result.Ok(), nil
}

func (a *DoneAction) GetVariants(_ adventuria.ActionContext) any {
	return nil
}

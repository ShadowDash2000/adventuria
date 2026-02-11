package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/result"
	"errors"
	"fmt"
)

type DropAction struct {
	adventuria.ActionBase
}

func (a *DropAction) CanDo(ctx adventuria.ActionContext) bool {
	currentCell, ok := ctx.User.CurrentCell()
	if ok {
		if currentCell.CantDrop() {
			return false
		}
	}

	if ctx.User.IsInJail() {
		return false
	}

	onBeforeDropCheckEvent := &adventuria.OnBeforeDropCheckEvent{
		AppContext:    ctx.AppContext,
		IsDropBlocked: false,
	}
	_, err := ctx.User.OnBeforeDropCheck().Trigger(onBeforeDropCheckEvent)
	if err != nil {
		return false
	}

	if onBeforeDropCheckEvent.IsDropBlocked {
		return false
	}

	return ctx.User.LastAction().Type() == ActionTypeRollWheel
}

func (a *DropAction) Do(ctx adventuria.ActionContext, req adventuria.ActionRequest) (*result.Result, error) {
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
			errors.New("drop.do(): current cell not found")
	}

	onBeforeDropEvent := &adventuria.OnBeforeDropEvent{
		AppContext:    ctx.AppContext,
		IsSafeDrop:    false,
		IsDropBlocked: false,
		PointsForDrop: adventuria.GameSettings.PointsForDrop(),
	}
	res, err := ctx.User.OnBeforeDrop().Trigger(onBeforeDropEvent)
	if err != nil {
		return res, err
	}
	if res.Failed() {
		return res, err
	}

	if onBeforeDropEvent.IsDropBlocked {
		return result.Err("drop is not allowed"), nil
	}

	action := ctx.User.LastAction()
	action.SetType(ActionTypeDrop)
	action.SetComment(comment)
	err = ctx.AppContext.App.Save(action.ProxyRecord())
	if err != nil {
		return result.Err("internal error: can't save action record"),
			fmt.Errorf("drop.do(): %w", err)
	}
	action.SetCanMove(true)

	if !onBeforeDropEvent.IsSafeDrop && !currentCell.IsSafeDrop() {
		ctx.User.SetPoints(ctx.User.Points() + onBeforeDropEvent.PointsForDrop)
		ctx.User.SetDropsInARow(ctx.User.DropsInARow() + 1)

		if !ctx.User.IsSafeDrop() {
			ctx.User.SetIsInJail(true)

			_, err = ctx.User.MoveToClosestCellType(ctx.AppContext, "jail")
			if err != nil {
				return result.Err("internal error: failed to move to the jail cell"),
					fmt.Errorf("drop.do(): %w", err)
			}
		}
	}

	res, err = ctx.User.OnAfterDrop().Trigger(&adventuria.OnAfterDropEvent{
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

func (a *DropAction) GetVariants(_ adventuria.ActionContext) any {
	return nil
}

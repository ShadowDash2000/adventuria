package actions

import (
	"adventuria/internal/adventuria"
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

func (a *DropAction) Do(ctx adventuria.ActionContext, req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
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
		}, errors.New("drop.do(): current cell not found")
	}

	onBeforeDropEvent := &adventuria.OnBeforeDropEvent{
		IsSafeDrop:    false,
		IsDropBlocked: false,
		PointsForDrop: adventuria.GameSettings.PointsForDrop(),
	}
	res, err := ctx.User.OnBeforeDrop().Trigger(onBeforeDropEvent)
	if res != nil && !res.Success {
		return &adventuria.ActionResult{
			Success: false,
			Error:   res.Error,
		}, err
	}
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: failed to trigger onBeforeDropEvent event",
		}, err
	}

	if onBeforeDropEvent.IsDropBlocked {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "drop is not allowed",
		}, nil
	}

	action := ctx.User.LastAction()
	action.SetType(ActionTypeDrop)
	action.SetComment(comment)
	err = adventuria.PocketBase.Save(action.ProxyRecord())
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't save action record",
		}, fmt.Errorf("drop.do(): %w", err)
	}
	action.SetCanMove(true)

	if !onBeforeDropEvent.IsSafeDrop && !currentCell.IsSafeDrop() {
		ctx.User.SetPoints(ctx.User.Points() + onBeforeDropEvent.PointsForDrop)
		ctx.User.SetDropsInARow(ctx.User.DropsInARow() + 1)

		if !ctx.User.IsSafeDrop() {
			ctx.User.SetIsInJail(true)

			_, err = ctx.User.MoveToClosestCellType("jail")
			if err != nil {
				return &adventuria.ActionResult{
					Success: false,
					Error:   "internal error",
				}, fmt.Errorf("drop.do(): %w", err)
			}
		}
	}

	res, err = ctx.User.OnAfterDrop().Trigger(&adventuria.OnAfterDropEvent{})
	if res != nil && !res.Success {
		return &adventuria.ActionResult{
			Success: false,
			Error:   res.Error,
		}, err
	}
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: failed to trigger onAfterDropEvent event",
		}, err
	}

	return &adventuria.ActionResult{
		Success: true,
	}, nil
}

func (a *DropAction) GetVariants(_ adventuria.ActionContext) any {
	return nil
}

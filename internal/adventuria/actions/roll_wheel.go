package actions

import (
	"adventuria/internal/adventuria"
	"errors"
	"fmt"
)

type RollWheelAction struct {
	adventuria.ActionBase
}

func (a *RollWheelAction) CanDo(ctx adventuria.ActionContext) bool {
	currentCell, ok := ctx.User.CurrentCell()
	if !ok {
		return false
	}

	if !currentCell.InCategory("activity") {
		return false
	}

	return !ctx.User.LastAction().CanMove() && ctx.User.LastAction().Type() != ActionTypeRollWheel
}

func (a *RollWheelAction) Do(ctx adventuria.ActionContext, req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	currentCell, ok := ctx.User.CurrentCell()
	if !ok {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: current cell not found",
		}, errors.New("roll_wheel.do(): current cell not found")
	}

	onBeforeWheelRollEvent := &adventuria.OnBeforeWheelRollEvent{
		AppContext:  ctx.AppContext,
		CurrentCell: currentCell.(adventuria.CellWheel),
	}
	eventRes, err := ctx.User.OnBeforeWheelRoll().Trigger(onBeforeWheelRollEvent)
	if eventRes != nil && !eventRes.Success {
		return &adventuria.ActionResult{
			Success: false,
			Error:   eventRes.Error,
		}, err
	}
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: failed to trigger onBeforeWheelRollEvent event",
		}, err
	}

	res, err := onBeforeWheelRollEvent.CurrentCell.Roll(ctx.AppContext, ctx.User, adventuria.RollWheelRequest(req))
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error",
		}, fmt.Errorf("roll_wheel.do(): %w", err)
	}

	action := ctx.User.LastAction()
	action.SetType(ActionTypeRollWheel)
	action.SetActivity(res.WinnerId)

	eventRes, err = ctx.User.OnAfterWheelRoll().Trigger(&adventuria.OnAfterWheelRollEvent{
		AppContext: ctx.AppContext,
		ItemId:     res.WinnerId,
	})
	if eventRes != nil && !eventRes.Success {
		return &adventuria.ActionResult{
			Success: false,
			Error:   eventRes.Error,
		}, err
	}
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: failed to trigger onAfterWheelRollEvent event",
		}, err
	}

	return &adventuria.ActionResult{
		Success: true,
		Data:    res,
	}, nil
}

func (a *RollWheelAction) GetVariants(_ adventuria.ActionContext) any {
	return nil
}

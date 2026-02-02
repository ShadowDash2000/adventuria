package actions

import (
	"adventuria/internal/adventuria"
	"fmt"
)

type RerollAction struct {
	adventuria.ActionBase
}

func (a *RerollAction) CanDo(ctx adventuria.ActionContext) bool {
	currentCell, ok := ctx.User.CurrentCell()
	if ok {
		if currentCell.CantReroll() {
			return false
		}
	}

	onBeforeRerollCheckEvent := &adventuria.OnBeforeRerollCheckEvent{
		IsRerollBlocked: false,
	}
	_, err := ctx.User.OnBeforeRerollCheck().Trigger(onBeforeRerollCheckEvent)
	if err != nil {
		return false
	}

	if onBeforeRerollCheckEvent.IsRerollBlocked {
		return false
	}

	return ctx.User.LastAction().Type() == ActionTypeRollWheel
}

func (a *RerollAction) Do(ctx adventuria.ActionContext, req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
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

	cell, ok := ctx.User.CurrentCell()
	if !ok {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: current cell not found",
		}, fmt.Errorf("reroll.do(): current cell not found")
	}

	cellWheel, ok := cell.(adventuria.CellWheel)
	if !ok {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: current cell isn't wheel cell",
		}, fmt.Errorf("reroll.do(): current cell isn't wheel cell")
	}

	action := ctx.User.LastAction()
	action.SetType(ActionTypeReroll)
	action.SetComment(comment)
	err := adventuria.PocketBase.Save(action.ProxyRecord())
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't save action record",
		}, fmt.Errorf("reroll.do(): %w", err)
	}
	action.MarkAsNew()

	err = cellWheel.RefreshItems(ctx.User)
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't refresh items on cell",
		}, fmt.Errorf("reroll.do(): %w", err)
	}

	res, err := ctx.User.OnAfterReroll().Trigger(&adventuria.OnAfterRerollEvent{})
	if res != nil && !res.Success {
		return &adventuria.ActionResult{
			Success: false,
			Error:   res.Error,
		}, err
	}
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: failed to trigger onAfterRerollEvent event",
		}, err
	}

	return &adventuria.ActionResult{
		Success: true,
	}, nil
}

package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/helper"
	"errors"
	"fmt"
)

type RollItemAction struct {
	adventuria.ActionBase
}

func (a *RollItemAction) CanDo(ctx adventuria.ActionContext) bool {
	if ctx.User.ItemWheelsCount() <= 0 {
		return false
	}

	if adventuria.GameActions.CanDo(ctx.AppContext, ctx.User, "done") {
		return false
	}

	return true
}

func (a *RollItemAction) Do(ctx adventuria.ActionContext, _ adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	res := &adventuria.WheelRollResult{}

	items := adventuria.GameItems.GetAllRollable()
	if len(items) == 0 {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: no items found",
		}, errors.New("roll_item.do(): no items found")
	}

	for _, item := range items {
		res.FillerItems = append(res.FillerItems, adventuria.WheelItem{
			Id:   item.ID(),
			Name: item.Name(),
			Icon: item.Icon(),
		})
	}

	res.WinnerId = helper.RandomItemFromSlice(items).ID()

	_, err := ctx.User.Inventory().MustAddItemById(ctx.AppContext, res.WinnerId)
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't add item to inventory",
		}, fmt.Errorf("roll_item.do(): %w", err)
	}

	ctx.User.SetItemWheelsCount(ctx.User.ItemWheelsCount() - 1)

	eventRes, err := ctx.User.OnAfterItemRoll().Trigger(&adventuria.OnAfterItemRollEvent{
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
			Error:   "internal error: failed to trigger onAfterItemRollEvent event",
		}, err
	}

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

func (a *RollItemAction) GetVariants(_ adventuria.ActionContext) any {
	return nil
}

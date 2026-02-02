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

	if adventuria.GameActions.CanDo(ctx.User, "done") {
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

	_, err := ctx.User.Inventory().MustAddItemById(res.WinnerId)
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't add item to inventory",
		}, fmt.Errorf("roll_item.do(): %w", err)
	}

	ctx.User.SetItemWheelsCount(ctx.User.ItemWheelsCount() - 1)

	onAfterItemRollEvent := &adventuria.OnAfterItemRollEvent{
		ItemId: res.WinnerId,
	}
	eventRes, err := ctx.User.OnAfterItemRoll().Trigger(onAfterItemRollEvent)
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

	onAfterWheelRollEvent := &adventuria.OnAfterWheelRollEvent{
		ItemId: res.WinnerId,
	}
	eventRes, err = ctx.User.OnAfterWheelRoll().Trigger(onAfterWheelRollEvent)
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

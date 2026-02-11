package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/helper"
	"adventuria/pkg/result"
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

func (a *RollItemAction) Do(ctx adventuria.ActionContext, _ adventuria.ActionRequest) (*result.Result, error) {
	res := &adventuria.WheelRollResult{}

	var items []adventuria.ItemRecord
	for item := range adventuria.GameItems.GetAllRollable() {
		items = append(items, item)
	}

	if len(items) == 0 {
		return result.Err("no items to roll"), nil
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
		return result.Err("internal error: failed to add item to the inventory"),
			fmt.Errorf("roll_item.do(): %w", err)
	}

	ctx.User.SetItemWheelsCount(ctx.User.ItemWheelsCount() - 1)

	eventRes, err := ctx.User.OnAfterItemRoll().Trigger(&adventuria.OnAfterItemRollEvent{
		AppContext: ctx.AppContext,
		ItemId:     res.WinnerId,
	})
	if err != nil {
		return eventRes, err
	}
	if eventRes.Failed() {
		return eventRes, err
	}

	eventRes, err = ctx.User.OnAfterWheelRoll().Trigger(&adventuria.OnAfterWheelRollEvent{
		AppContext: ctx.AppContext,
		ItemId:     res.WinnerId,
	})
	if err != nil {
		return eventRes, err
	}
	if eventRes.Failed() {
		return eventRes, err
	}

	return result.Ok().WithData(res), nil
}

func (a *RollItemAction) GetVariants(_ adventuria.ActionContext) any {
	return nil
}

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

func (a *RollItemAction) CanDo(user adventuria.User) bool {
	return user.ItemWheelsCount() > 0
}

func (a *RollItemAction) Do(user adventuria.User, _ adventuria.ActionRequest) (*adventuria.ActionResult, error) {
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

	_, err := user.Inventory().MustAddItemById(res.WinnerId)
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't add item to inventory",
		}, fmt.Errorf("roll_item.do(): %w", err)
	}

	user.SetItemWheelsCount(user.ItemWheelsCount() - 1)

	onAfterItemRollEvent := &adventuria.OnAfterItemRollEvent{
		ItemId: res.WinnerId,
	}
	err = user.OnAfterItemRoll().Trigger(onAfterItemRollEvent)
	if err != nil {
		adventuria.PocketBase.Logger().Error(
			"roll_item.do(): failed to trigger onAfterItemRoll event",
			"error",
			err,
		)
	}

	onAfterWheelRollEvent := &adventuria.OnAfterWheelRollEvent{
		ItemId: res.WinnerId,
	}
	err = user.OnAfterWheelRoll().Trigger(onAfterWheelRollEvent)
	if err != nil {
		adventuria.PocketBase.Logger().Error(
			"roll_item.do(): failed to trigger onAfterWheelRoll event",
			"error",
			err,
		)
	}

	return &adventuria.ActionResult{
		Success: true,
		Data:    res,
	}, nil
}

package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/helper"
	"errors"
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
		return nil, errors.New("items not found")
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
		return nil, err
	}

	user.SetItemWheelsCount(user.ItemWheelsCount() - 1)

	onAfterItemRollEvent := &adventuria.OnAfterItemRollEvent{
		ItemId: res.WinnerId,
	}
	err = user.OnAfterItemRoll().Trigger(onAfterItemRollEvent)
	if err != nil {
		return nil, err
	}

	onAfterWheelRollEvent := &adventuria.OnAfterWheelRollEvent{
		ItemId: res.WinnerId,
	}
	err = user.OnAfterWheelRoll().Trigger(onAfterWheelRollEvent)
	if err != nil {
		return nil, err
	}

	return &adventuria.ActionResult{
		Success: true,
		Data:    res,
	}, nil
}

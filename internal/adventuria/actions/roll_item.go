package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/helper"
	"errors"
)

type RollItemAction struct {
	adventuria.ActionBase
}

func (a *RollItemAction) CanDo() bool {
	return a.User().ItemWheelsCount() > 0
}

func (a *RollItemAction) NextAction() adventuria.ActionType {
	return ""
}

func (a *RollItemAction) Do(_ adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	itemsCol, err := adventuria.GameCollections.Get(adventuria.TableItems)
	if err != nil {
		return nil, err
	}

	res := &adventuria.WheelRollResult{
		Collection: itemsCol,
	}

	items := adventuria.GameItems.GetAllRollable()
	if len(items) == 0 {
		return nil, errors.New("items not found")
	}

	for _, item := range items {
		res.FillerItems = append(res.FillerItems, &adventuria.WheelItem{
			Name: item.Name(),
			Icon: item.Icon(),
		})
	}

	res.WinnerId = helper.RandomItemFromSlice(items).ID()

	_, err = a.User().Inventory().MustAddItemById(res.WinnerId)
	if err != nil {
		return nil, err
	}

	a.User().SetItemWheelsCount(a.User().ItemWheelsCount() - 1)

	onAfterItemRollEvent := &adventuria.OnAfterItemRollEvent{
		ItemId: res.WinnerId,
	}
	err = a.User().OnAfterItemRoll().Trigger(onAfterItemRollEvent)
	if err != nil {
		return nil, err
	}

	onAfterWheelRollEvent := &adventuria.OnAfterWheelRollEvent{
		ItemId: res.WinnerId,
	}
	err = a.User().OnAfterWheelRoll().Trigger(onAfterWheelRollEvent)
	if err != nil {
		return nil, err
	}

	return &adventuria.ActionResult{
		Success: true,
		Data:    res,
	}, nil
}

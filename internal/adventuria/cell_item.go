package adventuria

import (
	"adventuria/pkg/helper"
	"errors"
)

type CellItem struct {
	CellWheelBase
}

func NewCellItem() CellCreator {
	return func() Cell {
		return &CellItem{}
	}
}

func (c *CellItem) OnCellReached(user *User, _ *GameComponents) error {
	user.SetItemWheelsCount(user.ItemWheelsCount() + 1)
	return nil
}

func (c *CellItem) Roll(user *User) (*WheelRollResult, error) {
	itemsCol, err := c.gc.Cols.Get(TableItems)
	if err != nil {
		return nil, err
	}

	res := &WheelRollResult{
		Collection: itemsCol,
		EffectUse:  EffectUseOnRollItem,
	}

	items := c.gc.Items.GetAllRollable()

	if len(items) == 0 {
		return nil, errors.New("items not found")
	}

	for _, item := range items {
		res.FillerItems = append(res.FillerItems, &WheelItem{
			Name: item.Name(),
			Icon: item.Icon(),
		})
	}

	res.WinnerId = helper.RandomItemFromSlice(items).ID()

	err = user.Inventory.MustAddItemById(res.WinnerId)
	if err != nil {
		return nil, err
	}

	onAfterItemRollFields := &OnAfterItemRollFields{
		ItemId: res.WinnerId,
	}
	c.gc.Event.Go(OnAfterItemRoll, NewEventFields(user, c.gc, onAfterItemRollFields))

	return res, nil
}

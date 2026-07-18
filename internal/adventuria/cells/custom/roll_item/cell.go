package roll_item

import (
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/helper"
	"context"
	"errors"
)

type items interface {
	GetAllRollableByType(ctx context.Context, t model.ItemType) ([]*model.Item, error)
}

var _ model.Rollable = (*CellRollItem)(nil)

const Type model.CellType = "roll_item"

type CellRollItem struct {
	cells.CellBase
	items items
}

func NewDef(
	items items,
	categories ...string,
) cells.CellDef {
	return cells.NewCell(
		Type,
		func(cell model.CellInfo) model.Cell {
			return &CellRollItem{
				CellBase: cells.NewCellBase(cell),
				items:    items,
			}
		},
		categories...,
	)
}

func (c *CellRollItem) Roll(_ context.Context, _ *model.Events, player *model.Player) (*model.WheelRollResult, error) {
	itemsData := player.LastAction().DataList().Items

	if len(itemsData.Ids) == 0 {
		return nil, errors.New("no items to roll")
	}

	return &model.WheelRollResult{
		WinnerId: helper.RandomItemFromSlice(itemsData.Ids),
	}, nil
}

func (c *CellRollItem) RefreshItems(ctx context.Context, _ *model.Events, player *model.Player) error {
	return c.refreshItems(ctx, player)
}

func (c *CellRollItem) OnCellReached(ctx context.Context, _ *model.Events, player *model.Player, _ *model.ReachedContext) error {
	return c.refreshItems(ctx, player)
}

func (c *CellRollItem) OnCellLeft(_ context.Context, _ *model.Events, _ *model.Player) error {
	return nil
}

func (c *CellRollItem) refreshItems(ctx context.Context, player *model.Player) error {
	decodedValue, err := c.decodeValue(c.Value())
	if err != nil {
		return err
	}

	items, err := c.items.GetAllRollableByType(ctx, model.ItemType(decodedValue.ItemType))
	if err != nil {
		return err
	}

	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ID()
	}

	itemsData := player.LastAction().DataList().Items
	itemsData.Ids = ids
	player.LastAction().SetItemsData(itemsData)

	return nil
}

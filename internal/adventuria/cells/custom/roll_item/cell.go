package roll_item

import (
	"adventuria/internal/adventuria/board/actors/player"
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
var _ player.CellReached = (*CellRollItem)(nil)

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
	itemsState := player.LastAction().State().Items

	if len(itemsState.Ids) == 0 {
		return nil, errors.New("no items to roll")
	}

	return &model.WheelRollResult{
		WinnerId: helper.RandomItemFromSlice(itemsState.Ids),
	}, nil
}

func (c *CellRollItem) OnPlayerReached(ctx context.Context, actor *player.Actor, _ *model.ReachedContext) error {
	return c.refreshItems(ctx, actor.Player())
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

	player.LastAction().State().Items = &model.ActionItemsState{
		Ids: ids,
	}

	return nil
}

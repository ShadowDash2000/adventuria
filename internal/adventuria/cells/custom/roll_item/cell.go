package roll_item

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/helper"
	"context"
	"errors"
)

type itemsService interface {
	GetAllRollableByType(ctx context.Context, t model.ItemType) ([]*model.Item, error)
}

var _ cells.Rollable = (*CellRollItem)(nil)
var _ player.CellReached = (*CellRollItem)(nil)

const Type model.CellType = "roll_item"

type CellRollItem struct {
	cells.CellBase
	payload Payload
	items   itemsService
}

func NewDef(
	itemsService itemsService,
	categories ...string,
) cells.CellDef {
	return cells.NewCell(
		Type,
		func(cell model.CellInfo) (cells.Cell, error) {
			return New(itemsService, cell)
		},
		categories...,
	)
}

func New(itemsService itemsService, cell model.CellInfo) (*CellRollItem, error) {
	payload, err := NewPayloadFromJSON([]byte(cell.Value()))
	if err != nil {
		return nil, err
	}

	return &CellRollItem{
		CellBase: cells.NewCellBase(cell),
		payload:  payload,
		items:    itemsService,
	}, nil
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

func (c *CellRollItem) OnPlayerReached(ctx context.Context, actor *player.Actor, _ *cells.ReachedContext) error {
	return c.refreshItems(ctx, actor.Player())
}

func (c *CellRollItem) refreshItems(ctx context.Context, player *model.Player) error {
	items, err := c.items.GetAllRollableByType(ctx, model.ItemType(c.payload.ItemsType))
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

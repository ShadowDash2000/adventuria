package shop

import (
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
	"math/rand/v2"
)

type items interface {
	GetAllBuyableByType(ctx context.Context, t model.ItemType) ([]*model.Item, error)
}

var _ model.Refreshable = (*CellShop)(nil)

const shopMaxItems = 6

type CellShop struct {
	cells.CellBase
	itemsType model.ItemType
	items     items
}

func NewDef(
	cellType model.CellType,
	itemsType model.ItemType,
	items items,
	categories ...string,
) cells.CellDef {
	return cells.NewCell(
		cellType,
		func(cell model.CellInfo) model.Cell {
			return &CellShop{
				CellBase:  cells.NewCellBase(cell),
				itemsType: itemsType,
				items:     items,
			}
		},
		categories...,
	)
}

func (c *CellShop) OnCellReached(ctx context.Context, _ *model.Events, player *model.Player, _ *model.ReachedContext) error {
	err := c.refreshItems(ctx, player)
	if err != nil {
		return err
	}

	player.Progress().SetCanMove(true)

	return nil
}

func (c *CellShop) OnCellLeft(_ context.Context, _ *model.Events, _ *model.Player) error {
	return nil
}

func (c *CellShop) RefreshItems(ctx context.Context, _ *model.Events, player *model.Player) error {
	return c.refreshItems(ctx, player)
}

func (c *CellShop) refreshItems(ctx context.Context, player *model.Player) error {
	items, err := c.items.GetAllBuyableByType(ctx, c.itemsType)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return errors.New("no items to buy")
	}

	ids := make([]string, shopMaxItems)
	for i := range shopMaxItems {
		ids[i] = items[rand.N(len(items))].ID()
	}

	actionState := player.LastAction().State()
	actionState.Shop.Ids = ids
	player.LastAction().SetState(actionState)

	return nil
}

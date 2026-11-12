package shop

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"context"
	"math/rand/v2"
)

type itemsService interface {
	GetAllBuyableIDsByType(ctx context.Context, t model.ItemType) ([]string, error)
}

var _ player.CellReached = (*CellShop)(nil)

const shopMaxItems = 6

type CellShop struct {
	cells.CellBase
	itemsType model.ItemType
	items     itemsService
}

func NewDef(
	cellType model.CellType,
	itemsType model.ItemType,
	items itemsService,
	categories ...string,
) cells.CellDef {
	return cells.NewCell(
		cellType,
		func(cell model.CellInfo) (cells.Cell, error) {
			return &CellShop{
				CellBase:  cells.NewCellBase(cell),
				itemsType: itemsType,
				items:     items,
			}, nil
		},
		categories...,
	)
}

func (c *CellShop) OnPlayerReached(ctx context.Context, actor *player.Actor, _ *cells.ReachedContext) error {
	actor.Player().Progress().SetCanMove(true)

	ids, err := c.items.GetAllBuyableIDsByType(ctx, c.itemsType)
	if err != nil {
		return err
	}

	actionState := actor.Player().LastAction().State()
	actionState.ShopFilter = &model.ActionShopFilterState{
		ItemType: c.itemsType,
	}
	actionState.Shop, err = model.NewShopState(model.ActionShopStateCreate{
		Type: model.ShopTypeBuffet,
		Ids:  PickRandomIDs(ids),
	})
	if err != nil {
		return err
	}

	return nil
}

func PickRandomIDs(ids []string) []string {
	res := make([]string, 0, shopMaxItems)
	if len(ids) > 0 {
		for range shopMaxItems {
			res = append(res, ids[rand.N(len(ids))])
		}
	}
	return res
}

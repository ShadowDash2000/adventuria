package casino

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"context"
)

type items interface {
	GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error)
}

var _ player.CellReached = (*CellCasino)(nil)

const Type model.CellType = "casino"

type CellCasino struct {
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
			return &CellCasino{
				CellBase: cells.NewCellBase(cell),
				items:    items,
			}
		},
		categories...,
	)
}

func (c *CellCasino) OnPlayerReached(_ context.Context, actor *player.Actor, _ *model.ReachedContext) error {
	err := actor.Player().Progress().ItemWheelsCountChange(1)
	if err != nil {
		return err
	}

	actor.Player().Progress().SetCanMove(true)

	decodedValue, err := c.decodeValue(c.Value())
	if err != nil {
		return err
	}

	actionState := actor.Player().LastAction().State()
	actionState.ShopFilter = &model.ActionShopFilterState{
		Ids: decodedValue.ItemIds,
	}
	actionState.Shop, err = model.NewShopState(model.ActionShopStateCreate{
		Type:            model.ShopTypeCasino,
		Ids:             decodedValue.ItemIds,
		PriceMultiplier: decodedValue.PriceMultiplier,
	})
	if err != nil {
		return err
	}

	return nil
}

package casino

import (
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"context"
)

type items interface {
	GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error)
}

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

func (c *CellCasino) OnCellReached(_ context.Context, _ *model.Events, player *model.Player, _ *model.ReachedContext) error {
	err := player.Progress().ItemWheelsCountChange(1)
	if err != nil {
		return err
	}

	player.Progress().SetCanMove(true)

	decodedValue, err := c.decodeValue(c.Value())
	if err != nil {
		return err
	}

	actionState := player.LastAction().State()
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

func (c *CellCasino) OnCellLeft(_ context.Context, _ *model.Events, _ *model.Player) error {
	return nil
}

package casino

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"context"
)

type itemsService interface {
	GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error)
}

var _ player.CellReached = (*CellCasino)(nil)

const Type model.CellType = "casino"

type CellCasino struct {
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

func New(itemsService itemsService, cell model.CellInfo) (*CellCasino, error) {
	payload, err := NewPayloadFromJSON([]byte(cell.Value()))
	if err != nil {
		return nil, err
	}

	return &CellCasino{
		CellBase: cells.NewCellBase(cell),
		payload:  payload,
		items:    itemsService,
	}, nil
}

func (c *CellCasino) OnPlayerReached(_ context.Context, actor *player.Actor, _ *cells.ReachedContext) error {
	err := actor.Player().Progress().ItemWheelsCountChange(1)
	if err != nil {
		return err
	}

	actor.Player().Progress().SetCanMove(true)

	actionState := actor.Player().LastAction().State()
	actionState.ShopFilter = &model.ActionShopFilterState{
		Ids: c.payload.ItemIds,
	}
	actionState.Shop, err = model.NewShopState(model.ActionShopStateCreate{
		Type:            model.ShopTypeCasino,
		Ids:             c.payload.ItemIds,
		PriceMultiplier: c.payload.PriceMultiplier,
	})
	if err != nil {
		return err
	}

	return nil
}

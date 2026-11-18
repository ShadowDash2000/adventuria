package casino

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_actions/custom/shop"
	"context"
)

type itemsService interface {
	GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error)
}

type playerActions interface {
	Save(ctx context.Context, action *model.PlayerActionInfo) (*model.PlayerActionInfo, error)
}

var _ player.CellReached = (*CellCasino)(nil)

const Type model.CellType = "casino"

type CellCasino struct {
	cells.CellBase
	payload       Payload
	items         itemsService
	playerActions playerActions
}

func NewDef(
	itemsService itemsService,
	playerActions playerActions,
	categories ...string,
) cells.CellDef {
	return cells.NewCell(
		Type,
		func(cell model.CellInfo) (cells.Cell, error) {
			return New(itemsService, playerActions, cell)
		},
		categories...,
	)
}

func New(itemsService itemsService, playerActions playerActions, cell model.CellInfo) (*CellCasino, error) {
	payload, err := NewPayloadFromJSON([]byte(cell.Value()))
	if err != nil {
		return nil, err
	}

	return &CellCasino{
		CellBase:      cells.NewCellBase(cell),
		payload:       payload,
		items:         itemsService,
		playerActions: playerActions,
	}, nil
}

func (c *CellCasino) OnPlayerReached(ctx context.Context, actor *player.Actor, _ *cells.ReachedContext) error {
	err := actor.Player().Progress().ItemWheelsCountChange(1)
	if err != nil {
		return err
	}

	actor.Player().Progress().SetCanMove(true)

	shopAction, err := shop.New(shop.Payload{
		Type:            shop.TypeCasino,
		Ids:             c.payload.ItemIds,
		PriceMultiplier: c.payload.PriceMultiplier,
		Filter: shop.Filter{
			Ids: c.payload.ItemIds,
		},
	})
	if err != nil {
		return err
	}

	action, err := model.NewPlayerAction(model.PlayerActionCreate{
		PlayerId:     actor.Player().ID(),
		ActionId:     actor.Player().LastAction().ActionChainID(),
		PlayerAction: shopAction,
	})
	if err != nil {
		return err
	}

	_, err = c.playerActions.Save(ctx, action)
	if err != nil {
		return err
	}

	return nil
}

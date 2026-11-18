package shop

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_actions/custom/shop"
	"context"
)

type itemsService interface {
	GetAllBuyableIDsByType(ctx context.Context, t model.ItemType) ([]string, error)
}

type playerActions interface {
	Save(ctx context.Context, action *model.PlayerActionInfo) (*model.PlayerActionInfo, error)
}

var _ player.CellReached = (*CellShop)(nil)

const shopMaxItems = 6

type CellShop struct {
	cells.CellBase
	itemsType     model.ItemType
	items         itemsService
	playerActions playerActions
}

func NewDef(
	cellType model.CellType,
	itemsType model.ItemType,
	items itemsService,
	playerActions playerActions,
	categories ...string,
) cells.CellDef {
	return cells.NewCell(
		cellType,
		func(cell model.CellInfo) (cells.Cell, error) {
			return &CellShop{
				CellBase:      cells.NewCellBase(cell),
				itemsType:     itemsType,
				items:         items,
				playerActions: playerActions,
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

	shopAction, err := shop.New(shop.Payload{
		Type: shop.TypeBuffet,
		Filter: shop.Filter{
			ItemType: c.itemsType,
		},
	})
	if err != nil {
		return err
	}

	err = shopAction.Refresh(ids)
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

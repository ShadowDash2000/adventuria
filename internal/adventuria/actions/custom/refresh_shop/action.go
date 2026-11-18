package refresh_shop

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_actions/custom/shop"
	"context"
	"errors"
)

type itemsService interface {
	GetAllBuyableIDsByType(ctx context.Context, t model.ItemType) ([]string, error)
}

type playerActions interface {
	UpdatePayloadByID(ctx context.Context, id string, payload []byte) error
}

var _ actions.Action = (*RefreshShop)(nil)

const Type model.ActionType = "refresh_shop"

type RefreshShop struct {
	actions.ActionBase
	items         itemsService
	playerActions playerActions
}

func NewActionRefreshShopDef(items itemsService, playerActions playerActions) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() actions.Action {
			return &RefreshShop{
				ActionBase:    actions.NewActionBase(Type),
				items:         items,
				playerActions: playerActions,
			}
		},
	)
}

func (r *RefreshShop) CanDo(_ context.Context, actor *player.Actor) bool {
	return actor.Player().Actions().Has[*shop.Shop](shop.PlayerActionType)
}

func (r *RefreshShop) Do(ctx context.Context, actor *player.Actor, _ actions.Request) (any, error) {
	shopEntry, ok := actor.Player().Actions().Get[*shop.Shop](shop.PlayerActionType)
	if !ok {
		return nil, errs.ErrNoActiveShop
	}

	shopAction := shopEntry.Action()

	if actor.Player().Progress().Balance() < shopAction.RefreshPrice() {
		return nil, errs.ErrNotEnoughMoney
	}

	switch shopAction.ShopType() {
	case shop.TypeBuffet:
		err := r.doBuffetRefresh(ctx, shopAction)
		if err != nil {
			return nil, err
		}
	case shop.TypeCasino:
		err := r.doCasinoRefresh(shopAction)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unknown shop type")
	}

	err := actor.Player().Progress().BalanceChange(-shopAction.RefreshPrice())
	if err != nil {
		return nil, err
	}

	return nil, r.playerActions.UpdatePayloadByID(ctx, shopEntry.Data().ID(), shopAction.Payload())
}

package refresh_shop

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type itemsService interface {
	GetAllBuyableIDsByType(ctx context.Context, t model.ItemType) ([]string, error)
}

var _ actions.Action = (*RefreshShop)(nil)

const Type model.ActionType = "refresh_shop"

type RefreshShop struct {
	actions.ActionBase
	items itemsService
}

func NewActionRefreshShopDef(items itemsService) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() actions.Action {
			return &RefreshShop{
				ActionBase: actions.NewActionBase(Type),
				items:      items,
			}
		},
	)
}

func (r *RefreshShop) CanDo(_ context.Context, actor *player.Actor) bool {
	actionState := actor.Player().LastAction().State()
	if actionState.Shop == nil {
		return false
	}
	if actionState.ShopFilter == nil {
		return false
	}

	return true
}

func (r *RefreshShop) Do(ctx context.Context, actor *player.Actor, _ actions.Request) (any, error) {
	actionState := actor.Player().LastAction().State()
	if actionState.Shop == nil {
		return nil, errs.ErrNoActiveShop
	}
	if actionState.ShopFilter == nil {
		return nil, errs.ErrNoActiveShopFilter
	}

	if actor.Player().Progress().Balance() < actionState.Shop.RefreshPrice {
		return nil, errs.ErrNotEnoughMoney
	}

	switch actionState.Shop.Type {
	case model.ShopTypeBuffet:
		err := r.doBuffetRefresh(ctx, actor.Player())
		if err != nil {
			return nil, err
		}
	case model.ShopTypeCasino:
		r.doCasinoRefresh(actor.Player())
	default:
		return nil, errors.New("unknown shop type")
	}

	return nil, actor.Player().Progress().BalanceChange(-actionState.Shop.RefreshPrice)
}

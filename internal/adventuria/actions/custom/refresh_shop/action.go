package refresh_shop

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type itemsService interface {
	GetAllBuyableIDsByType(ctx context.Context, t model.ItemType) ([]string, error)
}

var _ model.Action = (*RefreshShop)(nil)

const Type model.ActionType = "refresh_shop"

type RefreshShop struct {
	actions.ActionBase
	items itemsService
}

func NewActionRefreshShopDef(items itemsService) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &RefreshShop{
				ActionBase: actions.NewActionBase(Type),
				items:      items,
			}
		},
	)
}

func (r *RefreshShop) CanDo(_ context.Context, _ *model.Events, player *model.Player) bool {
	actionState := player.LastAction().State()
	if actionState.Shop == nil {
		return false
	}
	if actionState.ShopFilter == nil {
		return false
	}

	return true
}

func (r *RefreshShop) Do(ctx context.Context, _ *model.Events, player *model.Player, _ model.ActionRequest) (any, error) {
	actionState := player.LastAction().State()
	if actionState.Shop == nil {
		return nil, errs.ErrNoActiveShop
	}
	if actionState.ShopFilter == nil {
		return nil, errs.ErrNoActiveShopFilter
	}

	if player.Progress().Balance() < actionState.Shop.RefreshPrice {
		return nil, errs.ErrNotEnoughMoney
	}

	switch actionState.Shop.Type {
	case model.ShopTypeBuffet:
		err := r.doBuffetRefresh(ctx, player, actionState)
		if err != nil {
			return nil, err
		}
	case model.ShopTypeCasino:
		r.doCasinoRefresh(player, actionState)
	default:
		return nil, errors.New("unknown shop type")
	}

	return nil, player.Progress().BalanceChange(-actionState.Shop.RefreshPrice)
}

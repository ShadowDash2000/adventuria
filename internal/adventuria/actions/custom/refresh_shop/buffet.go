package refresh_shop

import (
	"adventuria/internal/adventuria/cells/custom/shop"
	"adventuria/internal/adventuria/model"
	"context"
)

func (r *RefreshShop) doBuffetRefresh(ctx context.Context, player *model.Player) error {
	actionState := player.LastAction().State()
	ids, err := r.items.GetAllBuyableIDsByType(ctx, actionState.ShopFilter.ItemType)
	if err != nil {
		return err
	}

	actionState.Shop.Ids = shop.PickRandomIDs(ids)

	return nil
}

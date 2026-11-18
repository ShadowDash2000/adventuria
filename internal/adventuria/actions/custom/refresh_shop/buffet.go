package refresh_shop

import (
	"adventuria/internal/adventuria/player_actions/custom/shop"
	"context"
)

func (r *RefreshShop) doBuffetRefresh(ctx context.Context, shop *shop.Shop) error {
	ids, err := r.items.GetAllBuyableIDsByType(ctx, shop.Filter().ItemType)
	if err != nil {
		return err
	}

	return shop.Refresh(ids)
}

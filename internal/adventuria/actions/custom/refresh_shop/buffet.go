package refresh_shop

import (
	"adventuria/internal/adventuria/cells/custom/shop"
	"adventuria/internal/adventuria/model"
	"context"
)

func (r *RefreshShop) doBuffetRefresh(ctx context.Context, player *model.Player, state model.ActionState) error {
	ids, err := r.items.GetAllBuyableIDsByType(ctx, state.ShopFilter.ItemType)
	if err != nil {
		return err
	}

	state.Shop.Ids = shop.PickRandomIDs(ids)
	player.LastAction().SetState(state)

	return nil
}

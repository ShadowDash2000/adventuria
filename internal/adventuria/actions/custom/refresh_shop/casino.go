package refresh_shop

import (
	"adventuria/internal/adventuria/player_actions/custom/shop"
)

func (r *RefreshShop) doCasinoRefresh(shop *shop.Shop) error {
	return shop.SetItemIDs(shop.Filter().Ids)
}

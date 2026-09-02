package refresh_shop

import (
	"adventuria/internal/adventuria/model"
)

func (r *RefreshShop) doCasinoRefresh(player *model.Player) {
	actionState := player.LastAction().State()
	actionState.Shop.Ids = actionState.ShopFilter.Ids
}

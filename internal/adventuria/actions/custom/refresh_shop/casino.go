package refresh_shop

import (
	"adventuria/internal/adventuria/model"
)

func (r *RefreshShop) doCasinoRefresh(player *model.Player, state model.ActionState) {
	state.Shop.Ids = state.ShopFilter.Ids
	player.LastAction().SetState(state)
}

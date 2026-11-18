package refresh_shop

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/player_actions/custom/shop"
	"context"
)

var _ actions.View = (*RefreshShop)(nil)

func (r *RefreshShop) GetView(_ context.Context, actor *player.Actor) (any, error) {
	shopEntry, ok := actor.Player().Actions().Get[*shop.Shop](shop.PlayerActionType)
	if !ok {
		return nil, errs.ErrNoActiveShop
	}

	return struct {
		RefreshPrice int `json:"refresh_price"`
	}{
		RefreshPrice: shopEntry.Action().RefreshPrice(),
	}, nil
}

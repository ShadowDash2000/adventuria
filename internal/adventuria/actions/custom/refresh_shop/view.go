package refresh_shop

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/errs"
	"context"
)

var _ actions.View = (*RefreshShop)(nil)

func (r *RefreshShop) GetView(_ context.Context, actor *player.Actor) (any, error) {
	actionState := actor.Player().LastAction().State()
	if actionState.Shop == nil {
		return nil, errs.ErrNoActiveShop
	}

	return struct {
		RefreshPrice int `json:"refresh_price"`
	}{
		RefreshPrice: actionState.Shop.RefreshPrice,
	}, nil
}

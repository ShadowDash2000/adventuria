package refresh_shop

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
)

var _ model.WithView = (*RefreshShop)(nil)

func (r *RefreshShop) GetView(_ context.Context, _ *model.Events, player *model.Player) (any, error) {
	actionState := player.LastAction().State()
	if actionState.Shop == nil {
		return nil, errs.ErrNoActiveShop
	}

	return struct {
		RefreshPrice int `json:"refresh_price"`
	}{
		RefreshPrice: actionState.Shop.RefreshPrice,
	}, nil
}

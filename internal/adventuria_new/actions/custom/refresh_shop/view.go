package refresh_shop

import (
	"adventuria/internal/adventuria_new/model"
	"context"
)

var _ model.WithView = (*RefreshShop)(nil)

func (r *RefreshShop) GetView(ctx context.Context, _ *model.Events, player *model.Player) (any, error) {
	currentCell, err := r.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return nil, err
	}

	cellShopRefreshValue, err := r.decodeValue(currentCell.Data().Value())
	if err != nil {
		return nil, err
	}

	return struct {
		RefreshPrice int `json:"refresh_price"`
	}{
		RefreshPrice: cellShopRefreshValue.RefreshPrice,
	}, nil
}

package deal

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
)

var _ model.WithView = (*Deal)(nil)

func (d *Deal) GetView(_ context.Context, _ *model.Events, player *model.Player) (any, error) {
	deal := player.LastAction().State().Dealer

	if deal == nil {
		return nil, errs.ErrNoActiveDeals
	}

	return struct {
		Type model.DealType `json:"type"`
	}{
		Type: deal.Type,
	}, nil
}

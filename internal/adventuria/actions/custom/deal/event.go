package deal

import (
	"adventuria/internal/adventuria/model"
	"context"
)

var _ model.ActionEventCompatible = (*Deal)(nil)

func (d *Deal) CanDoOnEvent(_ context.Context, _ *model.Events, player *model.Player) bool {
	return player.LastAction().State().Dealer != nil
}

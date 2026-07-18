package buy

import (
	"adventuria/internal/adventuria/model"
	"context"
)

var _ model.ActionEventCompatible = (*Buy)(nil)

func (b *Buy) CanDoOnEvent(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

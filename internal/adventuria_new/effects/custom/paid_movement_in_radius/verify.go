package paid_movement_in_radius

import (
	"adventuria/internal/adventuria_new/model"
	"context"
)

var _ model.Verifiable = (*PaidMovementInRadius)(nil)

func (p *PaidMovementInRadius) Verify(_ context.Context, value string) error {
	_, err := p.decodeValue(value)
	return err
}

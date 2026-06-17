package points_change

import (
	"adventuria/internal/adventuria_new/model"
	"context"
)

var _ model.Verifiable = (*PointsChange)(nil)

func (p *PointsChange) Verify(_ context.Context, value string) error {
	_, err := p.decodeValue(value)
	return err
}

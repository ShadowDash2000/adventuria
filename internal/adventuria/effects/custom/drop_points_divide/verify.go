package drop_points_divide

import (
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

var _ model.Verifiable = (*DropPointsDivide)(nil)

func (d *DropPointsDivide) Verify(_ context.Context, value string) error {
	divider, err := d.decodeValue(value)
	if err != nil {
		return err
	}
	if divider <= 0 {
		return errors.New("divider must be greater than 0")
	}
	return nil
}

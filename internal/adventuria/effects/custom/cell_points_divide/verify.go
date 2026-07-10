package cell_points_divide

import (
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

var _ model.Verifiable = (*CellPointsDivide)(nil)

func (c *CellPointsDivide) Verify(_ context.Context, value string) error {
	divider, err := c.decodeValue(value)
	if err != nil {
		return err
	}
	if divider <= 0 {
		return errors.New("divider must be greater than 0")
	}
	return nil
}

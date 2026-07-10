package discount_price_divide

import (
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

var _ model.Verifiable = (*DiscountPriceDivide)(nil)

func (d *DiscountPriceDivide) Verify(_ context.Context, value string) error {
	divider, err := d.decodeValue(value)
	if err != nil {
		return err
	}
	if divider <= 0 {
		return errors.New("divider must be greater than 0")
	}
	return nil
}

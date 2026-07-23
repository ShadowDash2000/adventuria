package coins_for_item_dealer

import (
	"adventuria/internal/adventuria/model"
	"context"
)

var _ model.Verifiable = (*CoinsForItemDealer)(nil)

func (c *CoinsForItemDealer) Verify(ctx context.Context, value string) error {
	decodedValue, err := c.decodeValue(value)
	if err != nil {
		return err
	}

	_, err = c.items.GetByID(ctx, decodedValue.ItemId)
	if err != nil {
		return err
	}

	return nil
}

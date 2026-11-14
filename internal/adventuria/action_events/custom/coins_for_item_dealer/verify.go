package coins_for_item_dealer

import (
	"adventuria/internal/adventuria/model"
	"context"
)

var _ model.Verifiable = (*CoinsForItemDealer)(nil)

func (c *CoinsForItemDealer) Verify(ctx context.Context, value string) error {
	payload, err := NewPayloadFromJSON([]byte(value))
	if err != nil {
		return err
	}

	_, err = c.items.GetByID(ctx, payload.ItemId)
	if err != nil {
		return err
	}

	return nil
}

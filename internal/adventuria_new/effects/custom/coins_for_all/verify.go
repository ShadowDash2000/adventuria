package coins_for_all

import (
	"adventuria/internal/adventuria_new/model"
	"context"
)

var _ model.Verifiable = (*CoinsForAll)(nil)

func (c *CoinsForAll) Verify(_ context.Context, value string) error {
	_, err := c.decodeValue(value)
	if err != nil {
		return err
	}
	return nil
}

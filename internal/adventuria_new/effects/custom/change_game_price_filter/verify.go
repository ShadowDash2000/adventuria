package change_game_price_filter

import (
	"adventuria/internal/adventuria_new/model"
	"context"
)

var _ model.Verifiable = (*ChangeGamePriceFilter)(nil)

func (c *ChangeGamePriceFilter) Verify(_ context.Context, value string) error {
	_, err := c.decodeValue(value)
	if err != nil {
		return err
	}
	return nil
}

package change_dices

import (
	"adventuria/internal/adventuria/model"
	"context"
)

var _ model.Verifiable = (*ChangeDices)(nil)

func (c *ChangeDices) Verify(_ context.Context, value string) error {
	_, err := c.decodeValue(value)
	if err != nil {
		return err
	}
	return nil
}

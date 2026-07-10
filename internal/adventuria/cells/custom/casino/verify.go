package casino

import (
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
	"fmt"
)

var _ model.Verifiable = (*CellCasino)(nil)

func (c *CellCasino) Verify(ctx context.Context, value string) error {
	decodedValue, err := c.decodeValue(value)
	if err != nil {
		return err
	}

	if len(decodedValue.ItemIds) == 0 {
		return errors.New("item ids is empty")
	}

	items, err := c.items.GetByIDs(ctx, decodedValue.ItemIds)
	if err != nil {
		return err
	}

	if len(decodedValue.ItemIds) != len(items) {
		return fmt.Errorf("some of ids not found: %s", value)
	}

	return nil
}

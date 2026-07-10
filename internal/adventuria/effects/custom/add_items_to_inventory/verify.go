package add_items_to_inventory

import (
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
	"fmt"
)

var _ model.Verifiable = (*AddItemsToInventory)(nil)

func (a *AddItemsToInventory) Verify(ctx context.Context, value string) error {
	ids := a.decodeValue(value)
	if len(ids) == 0 {
		return errors.New("item ids is empty")
	}

	items, err := a.items.GetByIDs(ctx, ids)
	if err != nil {
		return err
	}

	if len(ids) != len(items) {
		return fmt.Errorf("some of ids not found: %s", value)
	}

	return nil
}

package roll_item

import (
	"adventuria/internal/adventuria_new/model"
	"context"
)

var _ model.Verifiable = (*CellRollItem)(nil)

func (c *CellRollItem) Verify(_ context.Context, value string) error {
	decodedValue, err := c.decodeValue(value)
	if err != nil {
		return err
	}

	_, err = model.NewItemType(decodedValue.ItemType)
	if err != nil {
		return err
	}

	return nil
}

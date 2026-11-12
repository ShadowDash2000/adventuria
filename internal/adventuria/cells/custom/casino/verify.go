package casino

import (
	"adventuria/internal/adventuria/model"
	"context"
	"fmt"
)

var _ model.Verifiable = (*CellCasino)(nil)

func (c *CellCasino) Verify(ctx context.Context, value string) error {
	payload, err := NewPayloadFromJSON([]byte(value))
	if err != nil {
		return err
	}

	items, err := c.items.GetByIDs(ctx, payload.ItemIds)
	if err != nil {
		return err
	}

	itemsMap := make(map[string]struct{}, len(items))
	for _, item := range items {
		itemsMap[item.ID()] = struct{}{}
	}

	for _, itemId := range payload.ItemIds {
		if _, ok := itemsMap[itemId]; !ok {
			return fmt.Errorf("item %s not found", itemId)
		}
	}

	return nil
}

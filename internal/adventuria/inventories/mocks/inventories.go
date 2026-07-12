package mocks

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type Inventories struct {
	AddItemByIDFunc func(ctx context.Context, events *model.Events, player *model.Player, itemId string) (*model.InventoryItem, error)
}

func (m *Inventories) AddItemByID(ctx context.Context, events *model.Events, player *model.Player, itemId string) (*model.InventoryItem, error) {
	if m.AddItemByIDFunc == nil {
		return nil, nil
	}

	return m.AddItemByIDFunc(ctx, events, player, itemId)
}

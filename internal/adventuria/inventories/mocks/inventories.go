package mocks

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/player_events/custom/item_received"
	"context"
)

type Inventories struct {
	AddItemByIDFunc func(ctx context.Context, actor *player.Actor, itemId string, source item_received.Source) (*effects.InventoryItemWithEffects, error)
}

func (m *Inventories) AddItemByID(ctx context.Context, actor *player.Actor, itemId string, source item_received.Source) (*effects.InventoryItemWithEffects, error) {
	if m.AddItemByIDFunc == nil {
		return nil, nil
	}

	return m.AddItemByIDFunc(ctx, actor, itemId, source)
}

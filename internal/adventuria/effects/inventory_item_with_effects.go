package effects

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/model"
	"context"
)

type InventoryItemWithEffects struct {
	model.InventoryItem
	effects []Effect
}

func RestoreInventoryItemWithEffects(inventory *model.Inventory, item *model.Item, effects []Effect) *InventoryItemWithEffects {
	return &InventoryItemWithEffects{
		InventoryItem: *model.RestoreInventoryItem(inventory, item),
		effects:       effects,
	}
}

func (i *InventoryItemWithEffects) Effects() []Effect {
	return i.effects
}

func (i *InventoryItemWithEffects) CanUse(ctx context.Context, actor *player.Actor) bool {
	if i.Inventory().IsActive() {
		return false
	}

	for _, effect := range i.effects {
		if !effect.CanUse(ctx, actor) {
			return false
		}
	}

	return true
}

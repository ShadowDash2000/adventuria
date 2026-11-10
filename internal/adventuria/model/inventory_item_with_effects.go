package model

import "context"

type InventoryItemWithEffects struct {
	InventoryItem
	effects []Effect
}

func RestoreInventoryItemWithEffects(inventory *Inventory, item *Item, effects []Effect) *InventoryItemWithEffects {
	return &InventoryItemWithEffects{
		inventory: inventory,
		item:      item,
		effects:   effects,
	}
}

func (i *InventoryItemWithEffects) Effects() []Effect {
	return i.effects
}

func (i *InventoryItemWithEffects) CanUse(ctx context.Context, events *Events, player *Player) bool {
	if i.inventory.IsActive() {
		return false
	}

	for _, effect := range i.effects {
		if !effect.CanUse(ctx, events, player) {
			return false
		}
	}

	return true
}

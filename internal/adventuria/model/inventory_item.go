package model

type InventoryItem struct {
	inventory *Inventory
	item      *Item
}

func RestoreInventoryItem(inventory *Inventory, item *Item) *InventoryItem {
	return &InventoryItem{inventory, item}
}

func (i *InventoryItem) Inventory() *Inventory {
	return i.inventory
}

func (i *InventoryItem) Item() *Item {
	return i.item
}

func (i *InventoryItem) UnappliedEffects() []string {
	appliedEffects := make(map[string]struct{}, len(i.inventory.AppliedEffects()))
	for _, effectId := range i.inventory.AppliedEffects() {
		appliedEffects[effectId] = struct{}{}
	}

	var unappliedEffects []string
	for _, effectId := range i.item.Effects() {
		if _, ok := appliedEffects[effectId]; !ok {
			unappliedEffects = append(unappliedEffects, effectId)
		}
	}

	return unappliedEffects
}

func (i *InventoryItem) CanDrop() bool {
	return i.item.CanDrop() && !i.inventory.IsActive()
}

package inventory

import "adventuria/internal/adventuria/model"

type playerInventoryView struct {
	PlayerId           string                     `json:"player_id"`
	PlayerName         string                     `json:"player_name"`
	Balance            int                        `json:"balance"`
	MaxInventorySlots  int                        `json:"max_inventory_slots"`
	UsedInventorySlots int                        `json:"used_inventory_slots"`
	Items              []*playerInventoryItemView `json:"items"`
}

type playerInventoryItemView struct {
	InventoryId      string                           `json:"inventory_id"`
	ItemId           string                           `json:"item_id"`
	Name             string                           `json:"name"`
	Icon             string                           `json:"icon"`
	Description      string                           `json:"description"`
	IsActive         bool                             `json:"is_active"`
	CanDrop          bool                             `json:"can_drop"`
	CanUse           bool                             `json:"can_use"`
	DropPrice        int                              `json:"drop_price"`
	UnappliedEffects []*playerInventoryItemEffectView `json:"unapplied_effects"`
}

type playerInventoryItemEffectView struct {
	Id   string           `json:"id"`
	Type model.EffectType `json:"type"`
}

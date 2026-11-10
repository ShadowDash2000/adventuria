package inventory

import "adventuria/internal/adventuria/inventories/model"

func playerInventoryToView(inventory *model.PlayerInventory) *playerInventoryView {
	return &playerInventoryView{
		PlayerId:           inventory.PlayerId,
		PlayerName:         inventory.PlayerName,
		Balance:            inventory.Balance,
		MaxInventorySlots:  inventory.MaxInventorySlots,
		UsedInventorySlots: inventory.UsedInventorySlots,
		Items:              playerInventoryItemsToView(inventory.Items),
	}
}

func playerInventoryItemToView(item *model.PlayerInventoryItem) *playerInventoryItemView {
	return &playerInventoryItemView{
		InventoryId:      item.InventoryId,
		ItemId:           item.ItemId,
		Name:             item.Name,
		Icon:             item.Icon,
		Description:      item.Description,
		IsActive:         item.IsActive,
		CanDrop:          item.CanDrop,
		CanUse:           item.CanUse,
		DropPrice:        item.DropPrice,
		UnappliedEffects: playerInventoryItemEffectsToView(item.UnappliedEffects),
	}
}

func playerInventoryItemsToView(items []*model.PlayerInventoryItem) []*playerInventoryItemView {
	res := make([]*playerInventoryItemView, len(items))
	for i, item := range items {
		res[i] = playerInventoryItemToView(item)
	}
	return res
}

func playerInventoryItemEffectToView(effect *model.PlayerInventoryItemEffect) *playerInventoryItemEffectView {
	return &playerInventoryItemEffectView{
		Id:   effect.Id,
		Type: effect.Type,
	}
}

func playerInventoryItemEffectsToView(effects []*model.PlayerInventoryItemEffect) []*playerInventoryItemEffectView {
	res := make([]*playerInventoryItemEffectView, len(effects))
	for i, effect := range effects {
		res[i] = playerInventoryItemEffectToView(effect)
	}
	return res
}

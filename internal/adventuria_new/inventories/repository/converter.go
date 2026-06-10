package repository

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/adventuria_new/model"

	"github.com/pocketbase/pocketbase/core"
)

func InventoryToRecord(inventory *model.Inventory, record *core.Record) {
	record.Id = inventory.ID()
	record.Set(schema.InventorySchema.Player, inventory.Player())
	record.Set(schema.InventorySchema.Item, inventory.Item())
	record.Set(schema.InventorySchema.IsActive, inventory.IsActive())
	record.Set(schema.InventorySchema.AppliedEffects, inventory.AppliedEffects())
}

func RecordToInventory(record *core.Record) *model.Inventory {
	return model.RestoreInventory(model.InventoryData{
		Id:             record.Id,
		Player:         record.GetString(schema.InventorySchema.Player),
		Item:           record.GetString(schema.InventorySchema.Item),
		IsActive:       record.GetBool(schema.InventorySchema.IsActive),
		AppliedEffects: record.GetStringSlice(schema.InventorySchema.AppliedEffects),
	})
}

func RecordToItem(record *core.Record) *model.Item {
	return model.RestoreItem(model.ItemData{
		Id:                record.Id,
		Name:              record.GetString(schema.ItemSchema.Name),
		Icon:              record.GetString(schema.ItemSchema.Icon),
		Effects:           record.GetStringSlice(schema.ItemSchema.Effects),
		IsUsingSlot:       record.GetBool(schema.ItemSchema.IsUsingSlot),
		IsActiveByDefault: record.GetBool(schema.ItemSchema.IsActiveByDefault),
		CanDrop:           record.GetBool(schema.ItemSchema.CanDrop),
		IsRollable:        record.GetBool(schema.ItemSchema.IsRollable),
		Description:       record.GetString(schema.ItemSchema.Description),
		Type:              model.ItemType(record.GetString(schema.ItemSchema.Type)),
		Price:             record.GetInt(schema.ItemSchema.Price),
	})
}

package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func InventoryToRecord(inventory *model.Inventory, record *core.Record) {
	record.Id = inventory.ID()
	record.Set(schema.InventorySchema.Activated, inventory.Activated())
	record.Set(schema.InventorySchema.Player, inventory.Player())
	record.Set(schema.InventorySchema.Item, inventory.Item())
	record.Set(schema.InventorySchema.IsActive, inventory.IsActive())
	record.Set(schema.InventorySchema.AppliedEffects, inventory.AppliedEffects())
}

func RecordToInventory(record *core.Record) *model.Inventory {
	return model.RestoreInventory(model.InventoryData{
		Id:             record.Id,
		Activated:      record.GetDateTime(schema.InventorySchema.Activated).Time(),
		Player:         record.GetString(schema.InventorySchema.Player),
		Item:           record.GetString(schema.InventorySchema.Item),
		IsActive:       record.GetBool(schema.InventorySchema.IsActive),
		AppliedEffects: record.GetStringSlice(schema.InventorySchema.AppliedEffects),
	})
}

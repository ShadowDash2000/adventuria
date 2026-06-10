package repository

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/adventuria_new/model"

	"github.com/pocketbase/pocketbase/core"
)

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

func RecordsToItems(records []*core.Record) []*model.Item {
	items := make([]*model.Item, len(records))
	for i, record := range records {
		items[i] = RecordToItem(record)
	}
	return items
}

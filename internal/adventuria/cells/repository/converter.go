package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToCell(record *core.Record) *model.CellInfo {
	return model.RestoreCellInfo(model.CellData{
		Id:                       record.Id,
		Disabled:                 record.GetBool(schema.CellSchema.Disabled),
		Sort:                     record.GetInt(schema.CellSchema.Sort),
		Type:                     model.CellType(record.GetString(schema.CellSchema.Type)),
		World:                    record.GetString(schema.CellSchema.World),
		Filter:                   record.GetString(schema.CellSchema.Filter),
		AudioPreset:              record.GetString(schema.CellSchema.AudioPreset),
		Icon:                     record.GetString(schema.CellSchema.Icon),
		Name:                     record.GetString(schema.CellSchema.Name),
		Points:                   record.GetInt(schema.CellSchema.Points),
		Coins:                    record.GetInt(schema.CellSchema.Coins),
		Description:              record.GetString(schema.CellSchema.Description),
		Color:                    record.GetString(schema.CellSchema.Color),
		CantDrop:                 record.GetBool(schema.CellSchema.CantDrop),
		CantReroll:               record.GetBool(schema.CellSchema.CantReroll),
		IsSafeDrop:               record.GetBool(schema.CellSchema.IsSafeDrop),
		IsCustomFilterNotAllowed: record.GetBool(schema.CellSchema.IsCustomFilterNotAllowed),
		IsChangeGameNotAllowed:   record.GetBool(schema.CellSchema.IsChangeGameNotAllowed),
		Value:                    record.GetString(schema.CellSchema.Value),
		LocalOrder:               record.GetInt("local_order"),
		GlobalOrder:              record.GetInt("global_order"),
	})
}

func RecordsToCells(records []*core.Record) []*model.CellInfo {
	cells := make([]*model.CellInfo, len(records))
	for i, record := range records {
		cells[i] = RecordToCell(record)
	}
	return cells
}

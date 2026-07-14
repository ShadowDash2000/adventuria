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
		EnergyConsume:            record.GetInt(schema.CellSchema.EnergyConsume),
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

func cellDTOToCell(dto *cellDTO) *model.CellInfo {
	value := ""
	if dto.Value != nil {
		value = *dto.Value
	}

	return model.RestoreCellInfo(model.CellData{
		Id:                       dto.Id,
		Disabled:                 dto.Disabled,
		Sort:                     dto.Sort,
		Type:                     model.CellType(dto.Type),
		World:                    dto.World,
		Filter:                   dto.Filter,
		AudioPreset:              dto.AudioPreset,
		Icon:                     dto.Icon,
		Name:                     dto.Name,
		Points:                   dto.Points,
		EnergyConsume:            dto.EnergyConsume,
		Coins:                    dto.Coins,
		Description:              dto.Description,
		Color:                    dto.Color,
		CantDrop:                 dto.CantDrop,
		CantReroll:               dto.CantReroll,
		IsSafeDrop:               dto.IsSafeDrop,
		IsCustomFilterNotAllowed: dto.IsCustomFilterNotAllowed,
		IsChangeGameNotAllowed:   dto.IsChangeGameNotAllowed,
		Value:                    value,
		LocalOrder:               dto.LocalOrder,
		GlobalOrder:              dto.GlobalOrder,
	})
}

func cellDTOsToCells(dtos []cellDTO) []*model.CellInfo {
	cells := make([]*model.CellInfo, len(dtos))
	for i := range dtos {
		cells[i] = cellDTOToCell(&dtos[i])
	}
	return cells
}

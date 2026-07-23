package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"strings"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToCellEventSchedule(record *core.Record) *model.CellEventSchedule {
	return model.RestoreCellEventSchedule(model.CellEventScheduleData{
		Id:              record.Id,
		ActionEvent:     record.GetString(schema.CellEventsScheduleSchema.ActionEvent),
		Effects:         record.GetStringSlice(schema.CellEventsScheduleSchema.Effects),
		ActiveCell:      record.GetString(schema.CellEventsScheduleSchema.ActiveCell),
		CellTypes:       stringToCellTypes(record.GetString(schema.CellEventsScheduleSchema.CellTypes)),
		Worlds:          record.GetStringSlice(schema.CellEventsScheduleSchema.Worlds),
		ShiftInterval:   record.GetInt(schema.CellEventsScheduleSchema.ShiftInterval),
		LastShiftChange: record.GetDateTime(schema.CellEventsScheduleSchema.LastShiftChange).Time(),
	})
}

func RecordsToCellEventSchedules(records []*core.Record) []*model.CellEventSchedule {
	res := make([]*model.CellEventSchedule, len(records))
	for i, record := range records {
		res[i] = RecordToCellEventSchedule(record)
	}
	return res
}

func stringToCellTypes(s string) []model.CellType {
	if s == "" {
		return nil
	}

	cellTypes := strings.Split(s, ";")
	res := make([]model.CellType, len(cellTypes))
	for i, cellType := range cellTypes {
		res[i] = model.CellType(cellType)
	}

	return res
}

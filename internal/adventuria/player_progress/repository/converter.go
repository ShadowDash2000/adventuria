package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func PlayerProgressToRecord(playerProgress *model.PlayerProgress, record *core.Record) {
	record.Id = playerProgress.ID()
	record.Set(schema.PlayerProgressSchema.Player, playerProgress.Player())
	record.Set(schema.PlayerProgressSchema.Season, playerProgress.Season())
	record.Set(schema.PlayerProgressSchema.CurrentWorld, playerProgress.CurrentWorld())
	record.Set(schema.PlayerProgressSchema.CanMove, playerProgress.CanMove())
	record.Set(schema.PlayerProgressSchema.Points, playerProgress.Points())
	record.Set(schema.PlayerProgressSchema.Balance, playerProgress.Balance())
	record.Set(schema.PlayerProgressSchema.Energy, playerProgress.Energy())
	record.Set(schema.PlayerProgressSchema.CellsPassed, playerProgress.CellsPassed())
	record.Set(schema.PlayerProgressSchema.IsInJail, playerProgress.IsInJail())
	record.Set(schema.PlayerProgressSchema.DropsInARow, playerProgress.DropsInARow())
	record.Set(schema.PlayerProgressSchema.ItemWheelsCount, playerProgress.ItemWheelsCount())
	record.Set(schema.PlayerProgressSchema.MaxInventorySlots, playerProgress.MaxInventorySlots())
}

func RecordToPlayerProgress(record *core.Record) *model.PlayerProgress {
	return model.RestorePlayerProgress(model.PlayerProgressData{
		Id:                record.Id,
		Player:            record.GetString(schema.PlayerProgressSchema.Player),
		Season:            record.GetString(schema.PlayerProgressSchema.Season),
		CurrentWorld:      record.GetString(schema.PlayerProgressSchema.CurrentWorld),
		CanMove:           record.GetBool(schema.PlayerProgressSchema.CanMove),
		Points:            record.GetInt(schema.PlayerProgressSchema.Points),
		Balance:           record.GetInt(schema.PlayerProgressSchema.Balance),
		Energy:            record.GetInt(schema.PlayerProgressSchema.Energy),
		CellsPassed:       record.GetInt(schema.PlayerProgressSchema.CellsPassed),
		IsInJail:          record.GetBool(schema.PlayerProgressSchema.IsInJail),
		DropsInARow:       record.GetInt(schema.PlayerProgressSchema.DropsInARow),
		ItemWheelsCount:   record.GetInt(schema.PlayerProgressSchema.ItemWheelsCount),
		MaxInventorySlots: record.GetInt(schema.PlayerProgressSchema.MaxInventorySlots),
	})
}

func RecordsToPlayerProgresses(records []*core.Record) []*model.PlayerProgress {
	progresses := make([]*model.PlayerProgress, len(records))
	for i, record := range records {
		progresses[i] = RecordToPlayerProgress(record)
	}
	return progresses
}

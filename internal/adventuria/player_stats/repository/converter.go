package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToPlayerStats(record *core.Record) *model.PlayerStats {
	return model.RestorePlayerStats(model.PlayerStatsData{
		Id:           record.Id,
		Player:       record.GetString(schema.PlayerStatsSchema.Player),
		Season:       record.GetString(schema.PlayerStatsSchema.Season),
		Drops:        record.GetInt(schema.PlayerStatsSchema.Drops),
		Rerolls:      record.GetInt(schema.PlayerStatsSchema.Rerolls),
		WasInJail:    record.GetInt(schema.PlayerStatsSchema.WasInJail),
		ItemsUsed:    record.GetInt(schema.PlayerStatsSchema.ItemsUsed),
		DiceRolls:    record.GetInt(schema.PlayerStatsSchema.DiceRolls),
		MaxDiceRoll:  record.GetInt(schema.PlayerStatsSchema.MaxDiceRoll),
		WheelsRolled: record.GetInt(schema.PlayerStatsSchema.WheelsRolled),
	})
}

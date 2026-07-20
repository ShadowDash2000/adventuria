package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToPlayerStats(record *core.Record) (*model.PlayerStats, error) {
	var activitiesStats model.ActivitiesStats
	err := record.UnmarshalJSONField(schema.PlayerStatsSchema.Activities, &activitiesStats)
	if err != nil {
		return nil, err
	}

	return model.RestorePlayerStats(model.PlayerStatsData{
		Id:              record.Id,
		Player:          record.GetString(schema.PlayerStatsSchema.Player),
		Season:          record.GetString(schema.PlayerStatsSchema.Season),
		ActivitiesStats: activitiesStats,
		Drops:           record.GetInt(schema.PlayerStatsSchema.Drops),
		Rerolls:         record.GetInt(schema.PlayerStatsSchema.Rerolls),
		WasInJail:       record.GetInt(schema.PlayerStatsSchema.WasInJail),
		ItemsUsed:       record.GetInt(schema.PlayerStatsSchema.ItemsUsed),
		DiceRolls:       record.GetInt(schema.PlayerStatsSchema.DiceRolls),
		MaxDiceRoll:     record.GetInt(schema.PlayerStatsSchema.MaxDiceRoll),
		WheelsRolled:    record.GetInt(schema.PlayerStatsSchema.WheelsRolled),
	}), nil
}

func PlayerStatsToRecord(stats *model.PlayerStats, record *core.Record) {
	record.Id = stats.ID()
	record.Set(schema.PlayerStatsSchema.Player, stats.Player())
	record.Set(schema.PlayerStatsSchema.Season, stats.Season())
	record.Set(schema.PlayerStatsSchema.Activities, stats.ActivitiesStats())
	record.Set(schema.PlayerStatsSchema.Drops, stats.Drops())
	record.Set(schema.PlayerStatsSchema.Rerolls, stats.Rerolls())
	record.Set(schema.PlayerStatsSchema.WasInJail, stats.WasInJail())
	record.Set(schema.PlayerStatsSchema.ItemsUsed, stats.ItemsUsed())
	record.Set(schema.PlayerStatsSchema.DiceRolls, stats.DiceRolls())
	record.Set(schema.PlayerStatsSchema.MaxDiceRoll, stats.MaxDiceRoll())
	record.Set(schema.PlayerStatsSchema.WheelsRolled, stats.WheelsRolled())
}

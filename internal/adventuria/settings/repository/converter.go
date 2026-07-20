package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func SettingsToRecord(settings *model.Settings, record *core.Record) {
	record.Id = settings.Id()
	record.Set(schema.SettingsSchema.EventEnded, settings.EventEnded())
	record.Set(schema.SettingsSchema.CurrentSeason, settings.CurrentSeason())
	record.Set(schema.SettingsSchema.CurrentWeek, settings.CurrentWeek())
	record.Set(schema.SettingsSchema.BlockAllActions, settings.BlockAllActions())
	record.Set(schema.SettingsSchema.EnergyDefault, settings.EnergyDefault())
	record.Set(schema.SettingsSchema.MaxInventorySlots, settings.MaxInventorySlots())
	record.Set(schema.SettingsSchema.PointsForDrop, settings.PointsForDrop())
	record.Set(schema.SettingsSchema.DropsToJail, settings.DropsToJail())
	record.Set(schema.SettingsSchema.IgdbFilterGameTypes, settings.IgdbFilterGameTypes())
	record.Set(schema.SettingsSchema.IgdbFilterPlatforms, settings.IgdbFilterPlatforms())
	record.Set(schema.SettingsSchema.IgdbFilterFirstReleaseDateMin, settings.IgdbFilterFirstReleaseDateMin())
	record.Set(schema.SettingsSchema.IgdbFilterFirstReleaseDateMax, settings.IgdbFilterFirstReleaseDateMax())
	record.Set(schema.SettingsSchema.IgdbGamesParsed, settings.IgdbGamesParsed())
	record.Set(schema.SettingsSchema.DisableIgdbParser, settings.DisableIgdbParser())
	record.Set(schema.SettingsSchema.DisableIgdbGamesParser, settings.DisableIgdbGamesParser())
	record.Set(schema.SettingsSchema.DisableSteamParser, settings.DisableSteamParser())
	record.Set(schema.SettingsSchema.DisableCheapsharkParser, settings.DisableCheapsharkParser())
	record.Set(schema.SettingsSchema.DisableHltbParser, settings.DisableHltbParser())
	record.Set(schema.SettingsSchema.DisableRefreshHltbTime, settings.DisableRefreshHltbTime())
	record.Set(schema.SettingsSchema.KillParser, settings.KillParser())
	record.Set(schema.SettingsSchema.IgdbForceUpdateGames, settings.IgdbForceUpdateGames())
}

func recordToGameTypeIdDb(record *core.Record) string {
	return record.GetString(schema.GameTypeSchema.IdDb)
}

func recordsToGameTypeIdsDb(records []*core.Record) []string {
	ids := make([]string, len(records))
	for i, record := range records {
		ids[i] = recordToGameTypeIdDb(record)
	}
	return ids
}

func recordToGamePlatformIdDb(record *core.Record) string {
	return record.GetString(schema.PlatformSchema.IdDb)
}

func recordsToGamePlatformIdsDb(records []*core.Record) []string {
	ids := make([]string, len(records))
	for i, record := range records {
		ids[i] = recordToGamePlatformIdDb(record)
	}
	return ids
}

func RecordToSettings(record *core.Record) *model.Settings {
	return model.RestoreSettings(model.SettingsData{
		Id:                record.Id,
		EventEnded:        record.GetBool(schema.SettingsSchema.EventEnded),
		CurrentSeason:     record.GetString(schema.SettingsSchema.CurrentSeason),
		CurrentWeek:       record.GetInt(schema.SettingsSchema.CurrentWeek),
		BlockAllActions:   record.GetBool(schema.SettingsSchema.BlockAllActions),
		EnergyDefault:     record.GetInt(schema.SettingsSchema.EnergyDefault),
		MaxInventorySlots: record.GetInt(schema.SettingsSchema.MaxInventorySlots),
		PointsForDrop:     record.GetInt(schema.SettingsSchema.PointsForDrop),
		DropsToJail:       record.GetInt(schema.SettingsSchema.DropsToJail),
		IgdbFilter: model.IgdbFilter{
			GameTypes:      recordsToGameTypeIdsDb(record.ExpandedAll(schema.SettingsSchema.IgdbFilterGameTypes)),
			Platforms:      recordsToGamePlatformIdsDb(record.ExpandedAll(schema.SettingsSchema.IgdbFilterPlatforms)),
			ReleaseDateMin: record.GetDateTime(schema.SettingsSchema.IgdbFilterFirstReleaseDateMin).Time(),
			ReleaseDateMax: record.GetDateTime(schema.SettingsSchema.IgdbFilterFirstReleaseDateMax).Time(),
		},
		IgdbGamesParsed:         uint(record.GetInt(schema.SettingsSchema.IgdbGamesParsed)),
		DisableIgdbParser:       record.GetBool(schema.SettingsSchema.DisableIgdbParser),
		DisableIgdbGamesParser:  record.GetBool(schema.SettingsSchema.DisableIgdbGamesParser),
		DisableSteamParser:      record.GetBool(schema.SettingsSchema.DisableSteamParser),
		DisableCheapsharkParser: record.GetBool(schema.SettingsSchema.DisableCheapsharkParser),
		DisableHltbParser:       record.GetBool(schema.SettingsSchema.DisableHltbParser),
		DisableRefreshHltbTime:  record.GetBool(schema.SettingsSchema.DisableRefreshHltbTime),
		KillParser:              record.GetBool(schema.SettingsSchema.KillParser),
		IgdbForceUpdateGames:    record.GetBool(schema.SettingsSchema.IgdbForceUpdateGames),
	})
}

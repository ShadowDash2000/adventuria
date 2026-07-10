package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func SettingsToRecord(settings *model.Settings, record *core.Record) {
	record.Id = settings.Id()
	record.Set(schema.SettingsSchema.EventEnded, settings.EventEnded)
	record.Set(schema.SettingsSchema.CurrentSeason, settings.CurrentSeason)
	record.Set(schema.SettingsSchema.CurrentWeek, settings.CurrentWeek)
	record.Set(schema.SettingsSchema.BlockAllActions, settings.BlockAllActions)
	record.Set(schema.SettingsSchema.MaxInventorySlots, settings.MaxInventorySlots)
	record.Set(schema.SettingsSchema.PointsForDrop, settings.PointsForDrop)
	record.Set(schema.SettingsSchema.DropsToJail, settings.DropsToJail)
	record.Set(schema.SettingsSchema.IgdbGamesParsed, settings.IgdbGamesParsed)
	record.Set(schema.SettingsSchema.DisableIgdbParser, settings.DisableIgdbParser)
	record.Set(schema.SettingsSchema.DisableSteamParser, settings.DisableSteamParser)
	record.Set(schema.SettingsSchema.DisableCheapsharkParser, settings.DisableCheapsharkParser)
	record.Set(schema.SettingsSchema.DisableHltbParser, settings.DisableHltbParser)
	record.Set(schema.SettingsSchema.DisableRefreshHltbTime, settings.DisableRefreshHltbTime)
	record.Set(schema.SettingsSchema.KillParser, settings.KillParser)
	record.Set(schema.SettingsSchema.IgdbForceUpdateGames, settings.IgdbForceUpdateGames)
}

func RecordToSettings(record *core.Record) *model.Settings {
	return model.RestoreSettings(model.SettingsData{
		Id:                      record.Id,
		EventEnded:              record.GetBool(schema.SettingsSchema.EventEnded),
		CurrentSeason:           record.GetString(schema.SettingsSchema.CurrentSeason),
		CurrentWeek:             record.GetInt(schema.SettingsSchema.CurrentWeek),
		BlockAllActions:         record.GetBool(schema.SettingsSchema.BlockAllActions),
		MaxInventorySlots:       record.GetInt(schema.SettingsSchema.MaxInventorySlots),
		PointsForDrop:           record.GetInt(schema.SettingsSchema.PointsForDrop),
		DropsToJail:             record.GetInt(schema.SettingsSchema.DropsToJail),
		IgdbGamesParsed:         uint(record.GetInt(schema.SettingsSchema.IgdbGamesParsed)),
		DisableIgdbParser:       record.GetBool(schema.SettingsSchema.DisableIgdbParser),
		DisableSteamParser:      record.GetBool(schema.SettingsSchema.DisableSteamParser),
		DisableCheapsharkParser: record.GetBool(schema.SettingsSchema.DisableCheapsharkParser),
		DisableHltbParser:       record.GetBool(schema.SettingsSchema.DisableHltbParser),
		DisableRefreshHltbTime:  record.GetBool(schema.SettingsSchema.DisableRefreshHltbTime),
		KillParser:              record.GetBool(schema.SettingsSchema.KillParser),
		IgdbForceUpdateGames:    record.GetBool(schema.SettingsSchema.IgdbForceUpdateGames),
	})
}

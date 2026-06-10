package repository

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/adventuria_new/model"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToActivity(record *core.Record) *model.Activity {
	return model.RestoreActivity(model.ActivityData{
		Id:               record.Id,
		IdDb:             record.GetString(schema.ActivitySchema.IdDb),
		Type:             model.ActivityType(record.GetString(schema.ActivitySchema.Type)),
		Name:             record.GetString(schema.ActivitySchema.Name),
		Slug:             record.GetString(schema.ActivitySchema.Slug),
		ReleaseDate:      record.GetDateTime(schema.ActivitySchema.ReleaseDate).Time(),
		Platforms:        record.GetStringSlice(schema.ActivitySchema.Platforms),
		Developers:       record.GetStringSlice(schema.ActivitySchema.Developers),
		Publishers:       record.GetStringSlice(schema.ActivitySchema.Publishers),
		Genres:           record.GetStringSlice(schema.ActivitySchema.Genres),
		Tags:             record.GetStringSlice(schema.ActivitySchema.Tags),
		Themes:           record.GetStringSlice(schema.ActivitySchema.Themes),
		GameType:         record.GetString(schema.ActivitySchema.GameType),
		SteamAppId:       uint(record.GetInt(schema.ActivitySchema.SteamAppId)),
		SteamAppPrice:    uint(record.GetInt(schema.ActivitySchema.SteamAppPrice)),
		HltbId:           uint(record.GetInt(schema.ActivitySchema.HltbId)),
		HltbCampaignTime: record.GetFloat(schema.ActivitySchema.HltbCampaignTime),
		Cover:            record.GetString(schema.ActivitySchema.Cover),
		CoverAlt:         record.GetString(schema.ActivitySchema.CoverAlt),
		Checksum:         record.GetString(schema.ActivitySchema.Checksum),
	})
}

func RecordsToActivities(records []*core.Record) []*model.Activity {
	activities := make([]*model.Activity, len(records))
	for i, record := range records {
		activities[i] = RecordToActivity(record)
	}
	return activities
}

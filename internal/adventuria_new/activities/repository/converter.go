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

func RecordToPlatform(record *core.Record) *model.Platform {
	return model.RestorePlatform(model.PlatformData{
		Id:       record.Id,
		IdDb:     record.GetString(schema.PlatformSchema.IdDb),
		Name:     record.GetString(schema.PlatformSchema.Name),
		Checksum: record.GetString(schema.PlatformSchema.Checksum),
	})
}

func RecordsToPlatforms(records []*core.Record) []*model.Platform {
	platforms := make([]*model.Platform, len(records))
	for i, record := range records {
		platforms[i] = RecordToPlatform(record)
	}
	return platforms
}

func RecordToDeveloper(record *core.Record) *model.Developer {
	return model.RestoreDeveloper(model.DeveloperData{
		Id:       record.Id,
		IdDb:     record.GetString(schema.CompanySchema.IdDb),
		Name:     record.GetString(schema.CompanySchema.Name),
		Checksum: record.GetString(schema.CompanySchema.Checksum),
	})
}

func RecordsToDevelopers(records []*core.Record) []*model.Developer {
	developers := make([]*model.Developer, len(records))
	for i, record := range records {
		developers[i] = RecordToDeveloper(record)
	}
	return developers
}

func RecordToPublisher(record *core.Record) *model.Publisher {
	return model.RestorePublisher(model.PublisherData{
		Id:       record.Id,
		IdDb:     record.GetString(schema.CompanySchema.IdDb),
		Name:     record.GetString(schema.CompanySchema.Name),
		Checksum: record.GetString(schema.CompanySchema.Checksum),
	})
}

func RecordsToPublishers(records []*core.Record) []*model.Publisher {
	publishers := make([]*model.Publisher, len(records))
	for i, record := range records {
		publishers[i] = RecordToPublisher(record)
	}
	return publishers
}

func RecordToGenre(record *core.Record) *model.Genre {
	return model.RestoreGenre(model.GenreData{
		Id:       record.Id,
		IdDb:     record.GetString(schema.GenreSchema.IdDb),
		Name:     record.GetString(schema.GenreSchema.Name),
		Checksum: record.GetString(schema.GenreSchema.Checksum),
	})
}

func RecordsToGenres(records []*core.Record) []*model.Genre {
	genres := make([]*model.Genre, len(records))
	for i, record := range records {
		genres[i] = RecordToGenre(record)
	}
	return genres
}

func RecordToTag(record *core.Record) *model.Tag {
	return model.RestoreTag(model.TagData{
		Id:       record.Id,
		IdDb:     record.GetString(schema.TagSchema.IdDb),
		Name:     record.GetString(schema.TagSchema.Name),
		Checksum: record.GetString(schema.TagSchema.Checksum),
	})
}

func RecordsToTags(records []*core.Record) []*model.Tag {
	tags := make([]*model.Tag, len(records))
	for i, record := range records {
		tags[i] = RecordToTag(record)
	}
	return tags
}

func RecordToTheme(record *core.Record) *model.Theme {
	return model.RestoreTheme(model.ThemeData{
		Id:       record.Id,
		IdDb:     record.GetString(schema.ThemeSchema.IdDb),
		Name:     record.GetString(schema.ThemeSchema.Name),
		Checksum: record.GetString(schema.ThemeSchema.Checksum),
	})
}

func RecordsToThemes(records []*core.Record) []*model.Theme {
	themes := make([]*model.Theme, len(records))
	for i, record := range records {
		themes[i] = RecordToTheme(record)
	}
	return themes
}

func RecordToActivityViewDetailed(record *core.Record) *model.ActivityViewDetailed {
	return model.RestoreActivityViewDetailed(
		RecordToActivity(record),
		RecordsToPlatforms(record.ExpandedAll(schema.ActivitySchema.Platforms)),
		RecordsToDevelopers(record.ExpandedAll(schema.ActivitySchema.Developers)),
		RecordsToPublishers(record.ExpandedAll(schema.ActivitySchema.Publishers)),
		RecordsToGenres(record.ExpandedAll(schema.ActivitySchema.Genres)),
		RecordsToTags(record.ExpandedAll(schema.ActivitySchema.Tags)),
		RecordsToThemes(record.ExpandedAll(schema.ActivitySchema.Themes)),
	)
}

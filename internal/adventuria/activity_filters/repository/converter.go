package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToActivityFilter(record *core.Record) *model.ActivityFilterInfo {
	return model.RestoreActivityFilter(model.ActivityFilterData{
		Id:   record.Id,
		Name: record.GetString(schema.ActivityFilterSchema.Name),
		ActivityFilter: model.ActivityFilter{
			Type:            model.ActivityType(record.GetString(schema.ActivityFilterSchema.Type)),
			Platforms:       record.GetStringSlice(schema.ActivityFilterSchema.Platforms),
			PlatformsStrict: record.GetBool(schema.ActivityFilterSchema.PlatformsStrict),
			GameTypes:       record.GetStringSlice(schema.ActivityFilterSchema.GameTypes),
			Developers:      record.GetStringSlice(schema.ActivityFilterSchema.Developers),
			Publishers:      record.GetStringSlice(schema.ActivityFilterSchema.Publishers),
			Genres:          record.GetStringSlice(schema.ActivityFilterSchema.Genres),
			Tags:            record.GetStringSlice(schema.ActivityFilterSchema.Tags),
			Themes:          record.GetStringSlice(schema.ActivityFilterSchema.Themes),
			MinPrice:        record.GetInt(schema.ActivityFilterSchema.MinPrice),
			MaxPrice:        record.GetInt(schema.ActivityFilterSchema.MaxPrice),
			ReleaseDateFrom: record.GetDateTime(schema.ActivityFilterSchema.ReleaseDateFrom).Time(),
			ReleaseDateTo:   record.GetDateTime(schema.ActivityFilterSchema.ReleaseDateTo).Time(),
			MinCampaignTime: record.GetFloat(schema.ActivityFilterSchema.MinCampaignTime),
			MaxCampaignTime: record.GetFloat(schema.ActivityFilterSchema.MaxCampaignTime),
			Activities:      record.GetStringSlice(schema.ActivityFilterSchema.Activities),
		},
	})
}

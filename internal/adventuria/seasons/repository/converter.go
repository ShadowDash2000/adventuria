package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

func SeasonToRecord(season *model.Season, record *core.Record) error {
	dateStart, err := types.ParseDateTime(season.SeasonDateStart())
	if err != nil {
		return err
	}
	dateEnd, err := types.ParseDateTime(season.SeasonDateEnd())
	if err != nil {
		return err
	}

	record.Id = season.ID()
	record.Set(schema.SeasonSchema.Name, season.Name())
	record.Set(schema.SeasonSchema.Slug, season.Slug())
	record.Set(schema.SeasonSchema.SeasonDateStart, dateStart)
	record.Set(schema.SeasonSchema.SeasonDateEnd, dateEnd)
	return nil
}

func RecordToSeason(record *core.Record) *model.Season {
	return model.RestoreSeason(model.SeasonData{
		Id:              record.Id,
		Name:            record.GetString(schema.SeasonSchema.Name),
		Slug:            record.GetString(schema.SeasonSchema.Slug),
		SeasonDateStart: record.GetDateTime(schema.SeasonSchema.SeasonDateStart).Time(),
		SeasonDateEnd:   record.GetDateTime(schema.SeasonSchema.SeasonDateEnd).Time(),
	})
}

package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToWorld(record *core.Record) *model.World {
	return model.RestoreWorld(model.WorldData{
		Id:                record.Id,
		Name:              record.GetString(schema.WorldSchema.Name),
		Slug:              record.GetString(schema.WorldSchema.Slug),
		Sort:              record.GetInt(schema.WorldSchema.Sort),
		IsLoop:            record.GetBool(schema.WorldSchema.IsLoop),
		IsDefaultWorld:    record.GetBool(schema.WorldSchema.IsDefaultWorld),
		TransitionToWorld: record.GetString(schema.WorldSchema.TransitionToWorld),
		Effects:           record.GetStringSlice(schema.WorldSchema.Effects),
	})
}

func RecordsToWorlds(records []*core.Record) []*model.World {
	worlds := make([]*model.World, len(records))
	for i, record := range records {
		worlds[i] = RecordToWorld(record)
	}
	return worlds
}

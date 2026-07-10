package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToPlatform(record *core.Record) *model.Platform {
	return model.RestorePlatform(model.PlatformData{
		Id:       record.Id,
		IdDb:     record.GetString(schema.PlatformSchema.IdDb),
		Name:     record.GetString(schema.PlatformSchema.Name),
		Checksum: record.GetString(schema.PlatformSchema.Checksum),
	})
}

func PlatformToRecord(platform *model.Platform, record *core.Record) {
	record.Id = platform.ID()
	record.Set(schema.PlatformSchema.IdDb, platform.IdDb())
	record.Set(schema.PlatformSchema.Name, platform.Name())
	record.Set(schema.PlatformSchema.Checksum, platform.Checksum())
}

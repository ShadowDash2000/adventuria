package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToTag(record *core.Record) *model.Tag {
	return model.RestoreTag(model.TagData{
		Id:       record.Id,
		IdDb:     record.GetString(schema.TagSchema.IdDb),
		Name:     record.GetString(schema.TagSchema.Name),
		Checksum: record.GetString(schema.TagSchema.Checksum),
	})
}

func TagToRecord(tag *model.Tag, record *core.Record) {
	record.Id = tag.ID()
	record.Set(schema.TagSchema.IdDb, tag.IdDb())
	record.Set(schema.TagSchema.Name, tag.Name())
	record.Set(schema.TagSchema.Checksum, tag.Checksum())
}

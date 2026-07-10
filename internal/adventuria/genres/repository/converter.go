package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToGenre(record *core.Record) *model.Genre {
	return model.RestoreGenre(model.GenreData{
		Id:       record.Id,
		IdDb:     record.GetString(schema.GenreSchema.IdDb),
		Name:     record.GetString(schema.GenreSchema.Name),
		Checksum: record.GetString(schema.GenreSchema.Checksum),
	})
}

func GenreToRecord(genre *model.Genre, record *core.Record) {
	record.Id = genre.ID()
	record.Set(schema.GenreSchema.IdDb, genre.IdDb())
	record.Set(schema.GenreSchema.Name, genre.Name())
	record.Set(schema.GenreSchema.Checksum, genre.Checksum())
}

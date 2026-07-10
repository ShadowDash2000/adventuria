package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToGameType(record *core.Record) *model.GameType {
	return model.RestoreGameType(model.GameTypeData{
		Id:       record.Id,
		IdDb:     record.GetString(schema.GameTypeSchema.IdDb),
		Name:     record.GetString(schema.GameTypeSchema.Name),
		Checksum: record.GetString(schema.GameTypeSchema.Checksum),
	})
}

func GameTypeToRecord(gameType *model.GameType, record *core.Record) {
	record.Id = gameType.ID()
	record.Set(schema.GameTypeSchema.IdDb, gameType.IdDb())
	record.Set(schema.GameTypeSchema.Name, gameType.Name())
	record.Set(schema.GameTypeSchema.Checksum, gameType.Checksum())
}

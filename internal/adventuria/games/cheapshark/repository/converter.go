package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToCheapShark(record *core.Record) *model.CheapShark {
	return model.RestoreCheapShark(model.CheapSharkData{
		Id:    record.Id,
		IdDb:  record.GetInt(schema.CheapSharkSchema.IdDb),
		Name:  record.GetString(schema.CheapSharkSchema.Name),
		Price: record.GetFloat(schema.CheapSharkSchema.Price),
	})
}

func CheapSharkToRecord(cheapShark *model.CheapShark, record *core.Record) {
	record.Id = cheapShark.ID()
	record.Set(schema.CheapSharkSchema.IdDb, cheapShark.IdDb())
	record.Set(schema.CheapSharkSchema.Name, cheapShark.Name())
	record.Set(schema.CheapSharkSchema.Price, cheapShark.Price())
}

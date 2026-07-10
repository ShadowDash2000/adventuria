package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToTheme(record *core.Record) *model.Theme {
	return model.RestoreTheme(model.ThemeData{
		Id:       record.Id,
		IdDb:     record.GetString(schema.ThemeSchema.IdDb),
		Name:     record.GetString(schema.ThemeSchema.Name),
		Checksum: record.GetString(schema.ThemeSchema.Checksum),
	})
}

func ThemeToRecord(theme *model.Theme, record *core.Record) {
	record.Id = theme.ID()
	record.Set(schema.ThemeSchema.IdDb, theme.IdDb())
	record.Set(schema.ThemeSchema.Name, theme.Name())
	record.Set(schema.ThemeSchema.Checksum, theme.Checksum())
}

package repository

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/adventuria_new/model"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToHowLongToBeat(record *core.Record) *model.HowLongToBeat {
	return model.RestoreHowLongToBeat(model.HowLongToBeatData{
		Id:       record.Id,
		IdDb:     record.GetInt(schema.HowLongToBeatSchema.IdDb),
		Name:     record.GetString(schema.HowLongToBeatSchema.Name),
		Year:     record.GetInt(schema.HowLongToBeatSchema.Year),
		Campaign: record.GetFloat(schema.HowLongToBeatSchema.Campaign),
	})
}

func HowLongToBeatToRecord(hltb *model.HowLongToBeat, record *core.Record) {
	record.Id = hltb.ID()
	record.Set(schema.HowLongToBeatSchema.IdDb, hltb.IdDb())
	record.Set(schema.HowLongToBeatSchema.Name, hltb.Name())
	record.Set(schema.HowLongToBeatSchema.Year, hltb.Year())
	record.Set(schema.HowLongToBeatSchema.Campaign, hltb.Campaign())
}

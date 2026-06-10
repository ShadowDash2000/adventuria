package repository

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/adventuria_new/model"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToEffect(record *core.Record) *model.EffectInfo {
	return model.RestoreEffectInfo(model.EffectData{
		Id:    record.Id,
		Name:  record.GetString(schema.EffectSchema.Name),
		Type:  model.EffectType(record.GetString(schema.EffectSchema.Type)),
		Value: record.GetString(schema.EffectSchema.Value),
	})
}

func RecordsToEffects(records []*core.Record) []*model.EffectInfo {
	effects := make([]*model.EffectInfo, len(records))
	for i, record := range records {
		effects[i] = RecordToEffect(record)
	}
	return effects
}

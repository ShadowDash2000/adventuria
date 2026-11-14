package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
)

func RecordToActionEvent(record *core.Record) *model.ActionEventInfo {
	return model.RestoreActionEvent(model.ActionEventData{
		Id:    record.Id,
		Name:  record.GetString(schema.ActionEventsSchema.Name),
		Type:  model.ActionEventType(record.GetString(schema.ActionEventsSchema.Type)),
		Value: record.GetString(schema.ActionEventsSchema.Value),
	})
}

func RecordsToActionEvents(records []*core.Record) []*model.ActionEventInfo {
	res := make([]*model.ActionEventInfo, len(records))
	for i, record := range records {
		res[i] = RecordToActionEvent(record)
	}
	return res
}

package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func RecordToActionEvent(record *core.Record) *model.ActionEventInfo {
	return model.RestoreActionEvent(model.ActionEventData{
		Id:         record.Id,
		Name:       record.GetString(schema.ActionEventsSchema.Name),
		Type:       model.ActionEventType(record.GetString(schema.ActionEventsSchema.Type)),
		ActionType: model.ActionType(record.GetString(schema.ActionEventsSchema.ActionType)),
		Value:      record.GetString(schema.ActionEventsSchema.Value),
	})
}

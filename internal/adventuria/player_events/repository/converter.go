package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func PlayerEventToRecord(playerEvent *model.PlayerEventInfo, record *core.Record) {
	record.Id = playerEvent.ID()
	record.Set(schema.PlayerEventsSchema.Player, playerEvent.Player())
	record.Set(schema.PlayerEventsSchema.Season, playerEvent.Season())
	record.Set(schema.PlayerEventsSchema.Type, playerEvent.Type())
	record.Set(schema.PlayerEventsSchema.Action, playerEvent.Action())
	record.Set(schema.PlayerEventsSchema.Payload, playerEvent.Payload())
}

func RecordToPlayerEvent(record *core.Record) *model.PlayerEventInfo {
	return model.RestorePlayerEvent(model.PlayerEventData{
		Id:      record.Id,
		Player:  record.GetString(schema.PlayerEventsSchema.Player),
		Season:  record.GetString(schema.PlayerEventsSchema.Season),
		Type:    model.PlayerEventType(record.GetString(schema.PlayerEventsSchema.Type)),
		Action:  record.GetString(schema.PlayerEventsSchema.Action),
		Payload: record.GetString(schema.PlayerEventsSchema.Payload),
	})
}

package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
	"ganigeorgiev.com/pocketbase/pocketbase/tools/types"
)

func RecordToPlayerAction(record *core.Record) *model.PlayerActionInfo {
	return model.RestorePlayerAction(model.PlayerActionData{
		Id:          record.Id,
		Player:      record.GetString(schema.PlayerActionsSchema.Player),
		Action:      record.GetString(schema.PlayerActionsSchema.Action),
		ActionEvent: record.GetString(schema.PlayerActionsSchema.ActionEvent),
		Type:        model.PlayerActionType(record.GetString(schema.PlayerActionsSchema.Type)),
		Payload:     record.Get(schema.PlayerActionsSchema.Payload).(types.JSONRaw),
	})
}

func RecordsToPlayerActions(records []*core.Record) []*model.PlayerActionInfo {
	actions := make([]*model.PlayerActionInfo, len(records))
	for i, record := range records {
		actions[i] = RecordToPlayerAction(record)
	}
	return actions
}

func PlayerActionToRecord(action *model.PlayerActionInfo, record *core.Record) {
	record.Id = action.ID()
	record.Set(schema.PlayerActionsSchema.Player, action.Player())
	record.Set(schema.PlayerActionsSchema.Action, action.Action())
	record.Set(schema.PlayerActionsSchema.ActionEvent, action.ActionEvent())
	record.Set(schema.PlayerActionsSchema.Type, action.Type())
	record.Set(schema.PlayerActionsSchema.Payload, action.Payload())
}

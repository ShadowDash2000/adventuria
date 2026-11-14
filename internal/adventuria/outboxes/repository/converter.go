package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
	"ganigeorgiev.com/pocketbase/pocketbase/tools/types"
)

func OutboxToRecord(outbox *model.OutboxInfo, record *core.Record) {
	record.Id = outbox.ID()
	record.Set(schema.OutboxSchema.Type, outbox.Type())
	record.Set(schema.OutboxSchema.Payload, outbox.Payload())
	record.Set(schema.OutboxSchema.Status, outbox.Status())
}

func RecordToOutbox(record *core.Record) *model.OutboxInfo {
	return model.RestoreOutbox(model.OutboxData{
		Id:      record.Id,
		Type:    model.OutboxType(record.GetString(schema.OutboxSchema.Type)),
		Payload: record.Get(schema.OutboxSchema.Payload).(types.JSONRaw),
		Status:  model.OutboxStatus(record.GetString(schema.OutboxSchema.Status)),
	})
}

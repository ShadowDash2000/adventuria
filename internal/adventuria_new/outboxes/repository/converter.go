package repository

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/adventuria_new/model"

	"github.com/pocketbase/pocketbase/core"
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
		Payload: record.GetString(schema.OutboxSchema.Payload),
		Status:  model.OutboxStatus(record.GetString(schema.OutboxSchema.Status)),
	})
}

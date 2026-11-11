package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
)

func PlayerTaskToRecord(task *model.PlayerTaskInfo, record *core.Record) {
	record.Id = task.ID()
	record.Set(schema.PlayerTasksSchema.Completed, task.Completed())
	record.Set(schema.PlayerTasksSchema.Progress, task.Progress())
	record.Set(schema.PlayerTasksSchema.Type, task.Type())
	record.Set(schema.PlayerTasksSchema.Trigger, task.Trigger())
	record.Set(schema.PlayerTasksSchema.Payload, task.Payload())
}

func RecordToPlayerTask(record *core.Record) *model.PlayerTaskInfo {
	return model.RestorePlayerTaskInfo(model.PlayerTaskData{
		Id:        record.Id,
		Completed: record.GetDateTime(schema.PlayerTasksSchema.Completed).Time(),
		Progress:  record.GetString(schema.PlayerTasksSchema.Progress),
		Type:      model.PlayerTaskType(record.GetString(schema.PlayerTasksSchema.Type)),
		Trigger:   model.PlayerTaskTrigger(record.GetString(schema.PlayerTasksSchema.Trigger)),
		Payload:   record.GetString(schema.PlayerTasksSchema.Payload),
	})
}

package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func ActionToRecord(action *model.ActionInfo, record *core.Record) error {
	state, err := actionStateToDTO(action.State())
	if err != nil {
		return err
	}

	record.Id = action.ID()
	record.Set(schema.ActionSchema.Player, action.Player())
	record.Set(schema.ActionSchema.Cell, action.Cell())
	record.Set(schema.ActionSchema.Status, string(action.Status()))
	record.Set(schema.ActionSchema.Activity, action.Activity())
	record.Set(schema.ActionSchema.Review, action.Review())
	record.Set(schema.ActionSchema.CellsPassed, action.CellsPassed())
	record.Set(schema.ActionSchema.State, state)
	record.Set(schema.ActionSchema.UsedItems, action.UsedItems())

	return nil
}

func RecordToAction(record *core.Record) (*model.ActionInfo, error) {
	var stateDTO actionState
	err := record.UnmarshalJSONField(schema.ActionSchema.State, &stateDTO)
	if err != nil {
		return nil, err
	}
	state, err := actionStateFromDTO(stateDTO)
	if err != nil {
		return nil, err
	}

	var usedItems []string
	err = record.UnmarshalJSONField(schema.ActionSchema.UsedItems, &usedItems)
	if err != nil {
		return nil, err
	}

	return model.RestoreAction(model.ActionData{
		Id:          record.Id,
		Player:      record.GetString(schema.ActionSchema.Player),
		Cell:        record.GetString(schema.ActionSchema.Cell),
		Status:      model.ActionStatus(record.GetString(schema.ActionSchema.Status)),
		Activity:    record.GetString(schema.ActionSchema.Activity),
		Review:      record.GetString(schema.ActionSchema.Review),
		CellsPassed: record.GetInt(schema.ActionSchema.CellsPassed),
		State:       state,
		UsedItems:   usedItems,
	}), nil
}

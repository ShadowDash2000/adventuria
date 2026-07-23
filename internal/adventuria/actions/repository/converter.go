package repository

import (
	"adventuria/internal/adventuria/actions/repository/dto"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func ActionToRecord(action *model.ActionInfo, record *core.Record) error {
	state, err := dto.ActionStateToDTO(action.State())
	if err != nil {
		return err
	}

	record.Id = action.ID()
	record.Set(schema.ActionSchema.Player, action.Player())
	record.Set(schema.ActionSchema.Cell, action.Cell())
	record.Set(schema.ActionSchema.Type, string(action.Type()))
	record.Set(schema.ActionSchema.Activity, action.Activity())
	record.Set(schema.ActionSchema.Review, action.Review())
	record.Set(schema.ActionSchema.CellsPassed, action.CellsPassed())
	record.Set(schema.ActionSchema.State, state)
	record.Set(schema.ActionSchema.UsedItems, action.UsedItems())
	record.Set(schema.ActionSchema.CustomActivityFilter, dto.CustomActivityFilterToDTO(action.CustomActivityFilter()))

	return nil
}

func RecordToAction(record *core.Record) (*model.ActionInfo, error) {
	var stateDTO dto.ActionState
	err := record.UnmarshalJSONField(schema.ActionSchema.State, &stateDTO)
	if err != nil {
		return nil, err
	}
	state, err := dto.ActionStateFromDTO(stateDTO)
	if err != nil {
		return nil, err
	}

	var usedItems []string
	err = record.UnmarshalJSONField(schema.ActionSchema.UsedItems, &usedItems)
	if err != nil {
		return nil, err
	}

	var filterDTO dto.CustomActivityFilter
	err = record.UnmarshalJSONField(schema.ActionSchema.CustomActivityFilter, &filterDTO)
	if err != nil {
		return nil, err
	}

	return model.RestoreAction(model.ActionData{
		Id:                   record.Id,
		Player:               record.GetString(schema.ActionSchema.Player),
		Cell:                 record.GetString(schema.ActionSchema.Cell),
		Type:                 model.ActionType(record.GetString(schema.ActionSchema.Type)),
		Activity:             record.GetString(schema.ActionSchema.Activity),
		Review:               record.GetString(schema.ActionSchema.Review),
		CellsPassed:          record.GetInt(schema.ActionSchema.CellsPassed),
		State:                state,
		UsedItems:            usedItems,
		CustomActivityFilter: dto.CustomActivityFilterFromDTO(filterDTO),
	}), nil
}

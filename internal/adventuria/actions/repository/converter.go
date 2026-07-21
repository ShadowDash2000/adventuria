package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"

	"github.com/pocketbase/pocketbase/core"
)

func ActionToRecord(action *model.ActionInfo, record *core.Record) {
	record.Id = action.ID()
	record.Set(schema.ActionSchema.Player, action.Player())
	record.Set(schema.ActionSchema.Cell, action.Cell())
	record.Set(schema.ActionSchema.Type, string(action.Type()))
	record.Set(schema.ActionSchema.Activity, action.Activity())
	record.Set(schema.ActionSchema.Review, action.Review())
	record.Set(schema.ActionSchema.CellsPassed, action.CellsPassed())
	record.Set(schema.ActionSchema.DataList, actionDataListToDTO(action.DataList()))
	record.Set(schema.ActionSchema.UsedItems, action.UsedItems())
	record.Set(schema.ActionSchema.CustomActivityFilter, customActivityFilterToDTO(action.CustomActivityFilter()))
}

func RecordToAction(record *core.Record) (*model.ActionInfo, error) {
	var dataListDTO actionDataListDTO
	err := record.UnmarshalJSONField(schema.ActionSchema.DataList, &dataListDTO)
	if err != nil {
		return nil, err
	}
	var usedItems []string
	err = record.UnmarshalJSONField(schema.ActionSchema.UsedItems, &usedItems)
	if err != nil {
		return nil, err
	}
	var filterDTO customActivityFilterDTO
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
		DataList:             actionDataListFromDTO(dataListDTO),
		UsedItems:            usedItems,
		CustomActivityFilter: customActivityFilterFromDTO(filterDTO),
	}), nil
}

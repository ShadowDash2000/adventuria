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
	record.Set(schema.ActionSchema.DataList, action.DataList())
	record.Set(schema.ActionSchema.UsedItems, action.UsedItems())
	record.Set(schema.ActionSchema.CustomActivityFilter, action.CustomActivityFilter())
}

func RecordToAction(record *core.Record) (*model.ActionInfo, error) {
	var actionDataList model.ActionDataList
	err := record.UnmarshalJSONField(schema.ActionSchema.DataList, &actionDataList)
	if err != nil {
		return nil, err
	}
	var usedItems []string
	err = record.UnmarshalJSONField(schema.ActionSchema.UsedItems, &usedItems)
	if err != nil {
		return nil, err
	}
	var customActivityFilter model.CustomActivityFilter
	err = record.UnmarshalJSONField(schema.ActionSchema.CustomActivityFilter, &customActivityFilter)
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
		DataList:             actionDataList,
		UsedItems:            usedItems,
		CustomActivityFilter: customActivityFilter,
	}), nil
}

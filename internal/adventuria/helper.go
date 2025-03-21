package adventuria

import (
	"adventuria/pkg/helper"
	"errors"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func RandomItem(gc *GameComponents) (*Item, error) {
	items, err := gc.app.FindRecordsByFilter(
		TableItems,
		"isRollable = true",
		"",
		0,
		0,
	)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errors.New("items not found")
	}

	return NewBaseItem(helper.RandomItemFromSlice(items)), nil
}

func RandomPresetItem(presetId string, gc *GameComponents) (*core.Record, error) {
	wheelItems, err := gc.app.FindRecordsByFilter(
		TableWheelItems,
		"presets.id = {:presetId}",
		"",
		0,
		0,
		dbx.Params{
			"presetId": presetId,
		},
	)
	if err != nil {
		return nil, err
	}

	if len(wheelItems) == 0 {
		return nil, errors.New("wheel items for preset not found")
	}

	return helper.RandomItemFromSlice(wheelItems), nil
}

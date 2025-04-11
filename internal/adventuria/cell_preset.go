package adventuria

import (
	"adventuria/pkg/helper"
	"errors"
	"github.com/pocketbase/dbx"
)

type CellPreset struct {
	CellWheelBase
}

func NewCellPreset() CellCreator {
	return func() Cell {
		return &CellPreset{}
	}
}

func (c *CellPreset) NextStep(user *User) string {
	nextStepType := ""

	switch user.LastAction.Type() {
	case ActionTypeRoll,
		ActionTypeReroll:
		nextStepType = ActionTypeRollWheelPreset
	case ActionTypeRollWheelPreset:
		nextStepType = ActionTypeDone
	case ActionTypeDone,
		ActionTypeDrop:
		nextStepType = ActionTypeRoll
	default:
		nextStepType = ActionTypeRollWheelPreset
	}

	return nextStepType
}

func (c *CellPreset) Roll(_ *User) (*WheelRollResult, error) {
	if c.Preset() == "" {
		return nil, errors.New("preset is not set")
	}

	wheelItemsCol, err := c.gc.Cols.Get(TableWheelItems)
	if err != nil {
		return nil, err
	}

	res := &WheelRollResult{
		Collection: wheelItemsCol,
	}

	items, err := c.gc.App.FindRecordsByFilter(
		TableWheelItems,
		"presets.id = {:presetId}",
		"",
		0,
		0,
		dbx.Params{
			"presetId": c.Preset(),
		},
	)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errors.New("wheel items for preset not found")
	}

	for _, item := range items {
		res.FillerItems = append(res.FillerItems, &WheelItem{
			Name: item.GetString("name"),
			Icon: item.GetString("icon"),
		})
	}

	res.WinnerId = helper.RandomItemFromSlice(items).Id

	return res, nil
}

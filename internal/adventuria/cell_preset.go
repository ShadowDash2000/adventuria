package adventuria

import (
	"adventuria/pkg/helper"
	"errors"

	"github.com/pocketbase/dbx"
)

type CellPreset struct {
	CellBase
}

func NewCellPreset() CellCreator {
	return func() Cell {
		return &CellPreset{}
	}
}

func (c *CellPreset) NextStep(user *User) string {
	nextStepType := ""

	switch user.LastAction.Type() {
	case ActionTypeRollDice,
		ActionTypeReroll:
		nextStepType = ActionTypeRollWheel
	case ActionTypeRollWheel:
		nextStepType = ActionTypeChooseResult
	case ActionTypeChooseResult,
		ActionTypeDrop:
		nextStepType = ActionTypeRollDice
	default:
		nextStepType = ActionTypeRollWheel
	}

	return nextStepType
}

func (c *CellPreset) Roll(_ *User) (*WheelRollResult, error) {
	if c.Preset() == "" {
		return nil, errors.New("preset is not set")
	}

	wheelItemsCol, err := GameCollections.Get(TableWheelItems)
	if err != nil {
		return nil, err
	}

	res := &WheelRollResult{
		Collection: wheelItemsCol,
	}

	items, err := PocketBase.FindRecordsByFilter(
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

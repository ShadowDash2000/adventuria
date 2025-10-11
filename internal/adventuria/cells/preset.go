package cells

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/helper"
	"errors"

	"github.com/pocketbase/dbx"
)

type CellPreset struct {
	adventuria.CellBase
}

func NewCellPreset() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellPreset{
			CellBase: adventuria.CellBase{},
		}
	}
}

func (c *CellPreset) NextStep(user adventuria.User) string {
	nextStepType := ""

	switch user.LastAction().Type() {
	case adventuria.ActionTypeRollDice,
		adventuria.ActionTypeReroll:
		nextStepType = adventuria.ActionTypeRollWheel
	case adventuria.ActionTypeRollWheel:
		nextStepType = adventuria.ActionTypeChooseResult
	case adventuria.ActionTypeChooseResult,
		adventuria.ActionTypeDrop:
		nextStepType = adventuria.ActionTypeRollDice
	default:
		nextStepType = adventuria.ActionTypeRollWheel
	}

	return nextStepType
}

func (c *CellPreset) Roll(_ adventuria.User) (*adventuria.WheelRollResult, error) {
	if c.Preset() == "" {
		return nil, errors.New("preset is not set")
	}

	wheelItemsCol, err := adventuria.GameCollections.Get(adventuria.TableWheelItems)
	if err != nil {
		return nil, err
	}

	res := &adventuria.WheelRollResult{
		Collection: wheelItemsCol,
	}

	items, err := adventuria.PocketBase.FindRecordsByFilter(
		adventuria.TableWheelItems,
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
		res.FillerItems = append(res.FillerItems, &adventuria.WheelItem{
			Name: item.GetString("name"),
			Icon: item.GetString("icon"),
		})
	}

	res.WinnerId = helper.RandomItemFromSlice(items).Id

	return res, nil
}

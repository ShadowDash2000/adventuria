package activities

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type DropResultCalculator struct{}

func NewDropResultCalculator() *DropResultCalculator {
	return &DropResultCalculator{}
}

func (d *DropResultCalculator) Calculate(
	ctx context.Context,
	events *model.Events,
	cell *model.CellInfo,
	mode model.EffectRunMode,
) (*model.ActivityDropResult, error) {
	onActivityDrop := &model.OnActivityDropEvent{
		Mode: mode,
		Result: model.NewActivityDropResult(model.ActivityDropResultData{
			IsBlocked:     cell.CantDrop(),
			IsSafe:        cell.IsSafeDrop(),
			EnergyConsume: 0,
		}),
	}
	err := events.OnActivityDrop().Trigger(ctx, onActivityDrop)
	if err != nil {
		return nil, err
	}

	return onActivityDrop.Result, nil
}

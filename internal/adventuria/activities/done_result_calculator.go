package activities

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type DoneResultCalculator struct{}

func NewDoneResultCalculator() *DoneResultCalculator {
	return &DoneResultCalculator{}
}

func (d *DoneResultCalculator) Calculate(
	ctx context.Context,
	events *model.Events,
	cell *model.CellInfo,
	mode model.EffectRunMode,
) (*model.ActivityDoneResult, error) {
	onActivityComplete := &model.OnActivityDoneEvent{
		Mode: mode,
		Result: model.NewActivityDoneResult(model.ActivityDoneResultData{
			Points:        cell.Points(),
			EnergyConsume: cell.EnergyConsume(),
			Coins:         cell.Coins(),
		}),
	}
	err := events.OnActivityDone().Trigger(ctx, onActivityComplete)
	if err != nil {
		return nil, err
	}

	return onActivityComplete.Result, nil
}

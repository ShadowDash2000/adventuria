package activities

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type CompletionResultCalculator struct{}

func NewCompletionResultCalculator() *CompletionResultCalculator {
	return &CompletionResultCalculator{}
}

func (c *CompletionResultCalculator) Calculate(
	ctx context.Context,
	events *model.Events,
	cell *model.CellInfo,
	mode model.EffectRunMode,
) (*model.ActivityCompletionResult, error) {
	onCellComplete := &model.OnActivityCompleteEvent{
		Mode: mode,
		Result: model.NewActivityCompletionResult(model.ActivityCompletionResultData{
			Points:        cell.Points(),
			EnergyConsume: cell.EnergyConsume(),
			Coins:         cell.Coins(),
		}),
	}
	err := events.OnActivityComplete().Trigger(ctx, onCellComplete)
	if err != nil {
		return nil, err
	}

	return onCellComplete.Result, nil
}

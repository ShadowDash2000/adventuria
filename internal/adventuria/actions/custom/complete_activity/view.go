package complete_activity

import (
	"adventuria/internal/adventuria/model"
	"context"
)

var _ model.WithView = (*CompleteActivity)(nil)

func (c *CompleteActivity) GetView(ctx context.Context, events *model.Events, player *model.Player) (any, error) {
	currentCell, err := c.cells.GetByPlayer(ctx, player)
	if err != nil {
		return nil, err
	}

	onDone := &model.OnDoneEvent{
		Mode: model.EffectRunModePreview,
		Result: model.NewDoneResult(model.DoneResultData{
			Points:        currentCell.Points(),
			EnergyConsume: currentCell.EnergyConsume(),
			Coins:         currentCell.Coins(),
		}),
	}
	err = events.OnDone().Trigger(ctx, onDone)
	if err != nil {
		return nil, err
	}

	return struct {
		DonePoints        int `json:"done_points"`
		DoneEnergyConsume int `json:"done_energy_consume"`
		DoneCoins         int `json:"done_coins"`
	}{
		DonePoints:        onDone.Result.Points(),
		DoneEnergyConsume: onDone.Result.EnergyConsume(),
		DoneCoins:         onDone.Result.Coins(),
	}, nil
}

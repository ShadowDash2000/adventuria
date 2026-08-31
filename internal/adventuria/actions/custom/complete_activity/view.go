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

	activityResult, err := c.activityResultCalculator.Calculate(ctx, events, currentCell, model.EffectRunModePreview)
	if err != nil {
		return nil, err
	}

	return struct {
		DonePoints        int `json:"done_points"`
		DoneEnergyConsume int `json:"done_energy_consume"`
		DoneCoins         int `json:"done_coins"`
	}{
		DonePoints:        activityResult.Points(),
		DoneEnergyConsume: activityResult.EnergyConsume(),
		DoneCoins:         activityResult.Coins(),
	}, nil
}

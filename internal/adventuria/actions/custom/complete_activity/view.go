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

	doneResult, err := c.doneResultCalculator.Calculate(ctx, events, currentCell, model.EffectRunModePreview)
	if err != nil {
		return nil, err
	}

	dropResult, err := c.dropResultCalculator.Calculate(ctx, events, currentCell, model.EffectRunModePreview)
	if err != nil {
		return nil, err
	}

	return struct {
		DonePoints        int  `json:"done_points"`
		DoneEnergyConsume int  `json:"done_energy_consume"`
		DoneCoins         int  `json:"done_coins"`
		IsSafeDrop        bool `json:"is_safe_drop"`
	}{
		DonePoints:        doneResult.Points(),
		DoneEnergyConsume: doneResult.EnergyConsume(),
		DoneCoins:         doneResult.Coins(),
		IsSafeDrop:        dropResult.IsSafe(),
	}, nil
}

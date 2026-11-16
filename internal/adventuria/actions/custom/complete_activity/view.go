package complete_activity

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/model"
	"context"
)

var _ actions.View = (*CompleteActivity)(nil)

func (c *CompleteActivity) GetView(ctx context.Context, actor *player.Actor) (any, error) {
	doneResult, err := c.doneResultCalculator.Calculate(ctx, actor, model.EffectRunModePreview)
	if err != nil {
		return nil, err
	}

	dropResult, err := c.dropResultCalculator.Calculate(ctx, actor, model.EffectRunModePreview)
	if err != nil {
		return nil, err
	}

	return struct {
		DonePoints        int  `json:"done_points"`
		DoneEnergyConsume int  `json:"done_energy_consume"`
		DoneCoins         int  `json:"done_coins"`
		DropEnergyConsume int  `json:"drop_energy_consume"`
		IsSafeDrop        bool `json:"is_safe_drop"`
	}{
		DonePoints:        doneResult.Points(),
		DoneEnergyConsume: doneResult.EnergyConsume(),
		DoneCoins:         doneResult.Coins(),
		DropEnergyConsume: dropResult.EnergyConsume(),
		IsSafeDrop:        dropResult.IsSafe(),
	}, nil
}

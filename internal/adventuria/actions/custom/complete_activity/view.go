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

	event := model.OnCompleteActivityView{
		CellPoints:        currentCell.Points(),
		CellEnergyConsume: currentCell.EnergyConsume(),
		CellCoins:         currentCell.Coins(),
	}
	err = events.OnCompleteActivityView().Trigger(ctx, &event)
	if err != nil {
		return nil, err
	}

	return struct {
		DonePoints        int `json:"done_points"`
		DoneEnergyConsume int `json:"done_energy_consume"`
		DoneCoins         int `json:"done_coins"`
	}{
		DonePoints:        event.CellPoints,
		DoneEnergyConsume: event.CellEnergyConsume,
		DoneCoins:         event.CellCoins,
	}, nil
}

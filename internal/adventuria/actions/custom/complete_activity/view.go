package complete_activity

import (
	"adventuria/internal/adventuria/model"
	"context"
)

var _ model.WithView = (*CompleteActivity)(nil)

func (c *CompleteActivity) GetView(ctx context.Context, _ *model.Events, player *model.Player) (any, error) {
	currentCell, err := c.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return nil, err
	}

	return struct {
		DoneEnergyConsume int `json:"done_energy_consume"`
	}{
		DoneEnergyConsume: currentCell.Data().EnergyConsume(),
	}, nil
}

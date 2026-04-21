package cells

import "adventuria/internal/adventuria"

type CellStart struct {
	adventuria.CellRecord
}

func (c *CellStart) OnCellReached(ctx *adventuria.CellReachedContext) error {
	ctx.Player.LastAction().SetCanMove(true)
	return nil
}

func (c *CellStart) OnCellLeft(_ *adventuria.CellLeftContext) error {
	return nil
}

func (c *CellStart) Verify(_ adventuria.AppContext, _ string) error {
	return nil
}

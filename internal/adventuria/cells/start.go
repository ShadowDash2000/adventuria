package cells

import "adventuria/internal/adventuria"

type CellStart struct {
	adventuria.CellRecord
}

func (c *CellStart) OnCellReached(ctx *adventuria.CellReachedContext) error {
	ctx.User.LastAction().SetCanMove(true)
	return nil
}

func (c *CellStart) OnCellLeft(_ *adventuria.CellLeftContext) error {
	return nil
}

func (c *CellStart) Verify(_ string) error {
	return nil
}

package cells

import "adventuria/internal/adventuria"

type CellStart struct {
	adventuria.CellBase
}

func NewCellStart() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellStart{
			adventuria.CellBase{},
		}
	}
}

func (c *CellStart) OnCellReached(ctx *adventuria.CellReachedContext) error {
	ctx.User.LastAction().SetCanMove(true)
	return nil
}

func (c *CellStart) Verify(_ string) error {
	return nil
}

func (c *CellStart) DecodeValue(_ string) (any, error) {
	return nil, nil
}

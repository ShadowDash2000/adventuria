package cells

import "adventuria/internal/adventuria"

type CellItem struct {
	adventuria.CellBase
}

func NewCellItem() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellItem{
			adventuria.CellBase{},
		}
	}
}

func (c *CellItem) OnCellReached(ctx *adventuria.CellReachedContext) error {
	ctx.User.SetItemWheelsCount(ctx.User.ItemWheelsCount() + 1)
	ctx.User.LastAction().SetCanMove(true)
	return nil
}

func (c *CellItem) Verify(_ string) error {
	return nil
}

func (c *CellItem) DecodeValue(_ string) (any, error) {
	return nil, nil
}

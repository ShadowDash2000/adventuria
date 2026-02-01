package cells

import (
	"adventuria/internal/adventuria"
)

type CellJail struct {
	CellGame
}

func NewCellJail() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellJail{
			CellGame: CellGame{
				CellRecord: adventuria.CellRecord{},
			},
		}
	}
}

func (c *CellJail) OnCellReached(ctx *adventuria.CellReachedContext) error {
	if ctx.User.IsInJail() {
		ctx.User.LastAction().SetCanMove(false)
		if err := c.refreshItems(ctx.User); err != nil {
			return err
		}

		_, err := ctx.User.OnAfterGoToJail().Trigger(&adventuria.OnAfterGoToJailEvent{})
		if err != nil {
			return err
		}
	} else {
		ctx.User.LastAction().SetCanMove(true)
	}
	return nil
}

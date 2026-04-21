package cells

import (
	"adventuria/internal/adventuria"
)

type CellJail struct {
	CellActivity
}

func NewCellJail() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellJail{
			CellActivity{
				activityType: adventuria.ActivityTypeGame,
			},
		}
	}
}

func (c *CellJail) OnCellReached(ctx *adventuria.CellReachedContext) error {
	if ctx.Player.Progress().IsInJail() {
		ctx.Player.LastAction().SetCanMove(false)
		if err := c.refreshItems(ctx.AppContext, ctx.Player); err != nil {
			return err
		}

		_, err := ctx.Player.OnAfterGoToJail().Trigger(&adventuria.OnAfterGoToJailEvent{})
		if err != nil {
			return err
		}
	} else {
		ctx.Player.LastAction().SetCanMove(true)
	}
	return nil
}

func (c *CellJail) OnCellLeft(ctx *adventuria.CellLeftContext) error {
	// If a player somehow left a jail, we need to free them
	if ctx.Player.Progress().IsInJail() {
		ctx.Player.Progress().SetIsInJail(false)
		ctx.Player.Progress().SetDropsInARow(0)
	}

	return nil
}

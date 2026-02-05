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
	if ctx.User.IsInJail() {
		ctx.User.LastAction().SetCanMove(false)
		if err := c.refreshItems(ctx.AppContext, ctx.User); err != nil {
			return err
		}

		_, err := ctx.User.OnAfterGoToJail().Trigger(&adventuria.OnAfterGoToJailEvent{})
		if err != nil {
			return err
		}
	} else {
		ctx.User.LastAction().SetCanMove(true)
	}
	return ctx.App.Save(ctx.User.LastAction().ProxyRecord())
}

func (c *CellJail) OnCellLeft(ctx *adventuria.CellLeftContext) error {
	// If a player somehow left a jail, we need to free them
	if ctx.User.IsInJail() {
		ctx.User.SetIsInJail(false)
		ctx.User.SetDropsInARow(0)
	}

	return ctx.App.Save(ctx.User.ProxyRecord())
}

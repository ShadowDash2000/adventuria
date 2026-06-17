package start

import (
	"adventuria/internal/adventuria_new/cells"
	"adventuria/internal/adventuria_new/model"
	"context"
)

var _ model.Cell = (*CellStart)(nil)

const Type model.CellType = "start"

type CellStart struct {
	cells.CellBase
}

func NewDef() cells.CellDef {
	return cells.NewCell(
		Type,
		func(cell model.CellInfo) model.Cell {
			return &CellStart{
				cells.NewCellBase(cell),
			}
		},
	)
}

func (c *CellStart) OnCellReached(_ context.Context, _ *model.Events, player *model.Player, _ *model.ReachedContext) error {
	player.LastAction().SetCanMove(true)
	return nil
}

func (c *CellStart) OnCellLeft(_ context.Context, _ *model.Events, _ *model.Player) error {
	return nil
}

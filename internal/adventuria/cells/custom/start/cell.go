package start

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"context"
)

var _ model.Cell = (*CellStart)(nil)
var _ player.CellReached = (*CellStart)(nil)

const Type model.CellType = "start"

type CellStart struct {
	cells.CellBase
}

func NewDef() cells.CellDef {
	return cells.NewCell(
		Type,
		func(cell model.CellInfo) (model.Cell, error) {
			return &CellStart{
				cells.NewCellBase(cell),
			}, nil
		},
	)
}

func (c *CellStart) OnPlayerReached(_ context.Context, actor *player.Actor, _ *model.ReachedContext) error {
	actor.Player().Progress().SetCanMove(true)
	return nil
}

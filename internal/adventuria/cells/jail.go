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
				CellWheel: &adventuria.CellWheelBase{
					CellBase: adventuria.CellBase{},
				},
			},
		}
	}
}

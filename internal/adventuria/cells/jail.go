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
				CellBase: adventuria.CellBase{},
			},
		}
	}
}

func (c *CellJail) OnCellReached(_ adventuria.User) error {
	return nil
}

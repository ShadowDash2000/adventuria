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

func (c *CellJail) OnCellReached(user adventuria.User) error {
	if !user.IsInJail() {
		user.LastAction().SetCanMove(true)
	}
	return nil
}

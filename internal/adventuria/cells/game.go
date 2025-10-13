package cells

import "adventuria/internal/adventuria"

type CellGame struct {
	adventuria.CellBase
}

func NewCellGame() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellGame{
			adventuria.CellBase{},
		}
	}
}

func (c *CellGame) Roll(_ adventuria.User) (*adventuria.WheelRollResult, error) {
	res := &adventuria.WheelRollResult{}

	// TODO
	panic("implement me")

	return res, nil
}

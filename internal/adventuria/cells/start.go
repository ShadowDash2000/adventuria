package cells

import "adventuria/internal/adventuria"

type CellStart struct {
	adventuria.CellBase
}

func NewCellStart() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellStart{
			adventuria.CellBase{},
		}
	}
}

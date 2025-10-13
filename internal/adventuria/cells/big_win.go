package cells

import (
	"adventuria/internal/adventuria"
)

type CellBigWin struct {
	CellPreset
}

func NewCellBigWin() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellBigWin{
			CellPreset: CellPreset{
				CellBase: adventuria.CellBase{},
			},
		}
	}
}

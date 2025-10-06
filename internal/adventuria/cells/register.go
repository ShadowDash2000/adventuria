package cells

import "adventuria/internal/adventuria"

const (
	CellTypeGame   adventuria.CellType = "game"
	CellTypeStart  adventuria.CellType = "start"
	CellTypeJail   adventuria.CellType = "jail"
	CellTypePreset adventuria.CellType = "preset"
	CellTypeItem   adventuria.CellType = "item"
	CellTypeBigWin adventuria.CellType = "bigWin"
)

func WithBaseCells() {
	adventuria.RegisterCells(map[adventuria.CellType]adventuria.CellCreator{
		CellTypeGame:   NewCellGame(),
		CellTypeStart:  NewCellStart(),
		CellTypeJail:   NewCellJail(),
		CellTypePreset: NewCellPreset(),
		CellTypeItem:   NewCellItem(),
		CellTypeBigWin: NewCellBigWin(),
	})
}

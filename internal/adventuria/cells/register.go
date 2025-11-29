package cells

import "adventuria/internal/adventuria"

const (
	CellTypeGame     adventuria.CellType = "game"
	CellTypeStart    adventuria.CellType = "start"
	CellTypeJail     adventuria.CellType = "jail"
	CellTypeItem     adventuria.CellType = "item"
	CellTypeShop     adventuria.CellType = "shop"
	CellTypeTeleport adventuria.CellType = "teleport"
)

func WithBaseCells() {
	adventuria.RegisterCells(map[adventuria.CellType]adventuria.CellCreator{
		CellTypeGame:     NewCellGame(),
		CellTypeStart:    NewCellStart(),
		CellTypeJail:     NewCellJail(),
		CellTypeItem:     NewCellItem(),
		CellTypeShop:     NewCellShop(),
		CellTypeTeleport: NewCellTeleport(),
	})
}

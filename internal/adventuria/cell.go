package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
	"maps"
)

type CellType string

const (
	CellTypeGame   CellType = "game"
	CellTypeStart  CellType = "start"
	CellTypeJail   CellType = "jail"
	CellTypePreset CellType = "preset"
	CellTypeItem   CellType = "item"
)

type Cell interface {
	core.RecordProxy
	ID() string
	Sort() int
	Type() CellType
	SetType(CellType)
	Preset() string
	AudioPresets() []string
	Icon() string
	Name() string
	Points() int
	Description() string
	Color() string
	CantDrop() bool
	CantReroll() bool
	CantChooseAfterDrop() bool
	IsSafeDrop() bool
	NextStep(*User) string
	OnCellReached(*User) error
}

var CellsList = map[CellType]CellCreator{
	CellTypeGame:   NewCellGame(),
	CellTypeStart:  NewCellStart(),
	CellTypeJail:   NewCellJail(),
	CellTypePreset: NewCellPreset(),
	CellTypeItem:   NewCellItem(),
}

type CellCreator func() Cell

func RegisterCells(cells map[CellType]CellCreator) {
	maps.Insert(CellsList, maps.All(cells))
}

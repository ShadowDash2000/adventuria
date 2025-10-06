package adventuria

import (
	"maps"

	"github.com/pocketbase/pocketbase/core"
)

type CellType string

type Cell interface {
	core.RecordProxy
	ID() string
	IsActive() bool
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
	IsSafeDrop() bool
	NextStep(User) string
	OnCellReached(User) error
}

var CellsList = map[CellType]CellCreator{}

type CellCreator func(ServiceLocator) Cell

func RegisterCells(cells map[CellType]CellCreator) {
	maps.Insert(CellsList, maps.All(cells))
}

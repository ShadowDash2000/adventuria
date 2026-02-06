package adventuria

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

type CellType string

type Cell interface {
	core.RecordProxy
	ID() string
	Sort() int
	Type() CellType
	setType(CellType)
	Categories() []string
	InCategory(string) bool
	Filter() string
	AudioPreset() []string
	Icon() string
	Name() string
	Points() int
	Coins() int
	Description() string
	Color() string
	CantDrop() bool
	CantReroll() bool
	IsSafeDrop() bool
	IsCustomFilterNotAllowed() bool
	OnCellReached(*CellReachedContext) error
	OnCellLeft(*CellLeftContext) error
	Verify(AppContext, string) error
	Value() string
	UnmarshalValue(result any) error
}

type CellReachedContext struct {
	AppContext
	User  User
	Moves []*MoveResult
}

type CellLeftContext struct {
	AppContext
	User User
}

var cellsList = map[CellType]CellDef{}

type CellCreator func() Cell

type CellDef struct {
	Type       CellType
	Categories []string
	New        func(record *core.Record) Cell
}

func RegisterCells(cells []CellDef) {
	for _, cellDef := range cells {
		cellsList[cellDef.Type] = cellDef
	}
}

func IsCellTypeExist(t CellType) bool {
	_, ok := cellsList[t]
	return ok
}

func NewCell(t CellType, newCellFn CellCreator, categories ...string) CellDef {
	return CellDef{
		Type:       t,
		Categories: categories,
		New: func(record *core.Record) Cell {
			c := newCellFn()
			c.setType(t)
			return c
		},
	}
}

func NewCellFromRecord(record *core.Record) (Cell, error) {
	t := CellType(record.GetString("type"))

	cellDef, ok := cellsList[t]
	if !ok {
		return nil, fmt.Errorf("unknown cell type: %s", t)
	}

	cell := cellDef.New(record)
	cell.SetProxyRecord(record)

	return cell, nil
}

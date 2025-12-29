package adventuria

import (
	"fmt"
	"maps"

	"github.com/pocketbase/pocketbase/core"
)

type CellType string

type Cell interface {
	core.RecordProxy
	ID() string
	Sort() int
	Type() CellType
	SetType(CellType)
	Filter() string
	AudioPresets() []string
	Icon() string
	Name() string
	Points() int
	Coins() int
	Description() string
	Color() string
	CantDrop() bool
	CantReroll() bool
	IsSafeDrop() bool
	OnCellReached(*CellReachedContext) error
	Verify(string) error
	DecodeValue(string) (any, error)
	Value() string
}

var cellsList = map[CellType]CellCreator{}

type CellCreator func() Cell

func RegisterCells(cells map[CellType]CellCreator) {
	maps.Insert(cellsList, maps.All(cells))
}

func IsCellTypeExist(t CellType) bool {
	_, ok := cellsList[t]
	return ok
}

func NewCellFromRecord(record *core.Record) (Cell, error) {
	t := CellType(record.GetString("type"))

	cellCreator, ok := cellsList[t]
	if !ok {
		return nil, fmt.Errorf("unknown cell type: %s", t)
	}

	cell := cellCreator()
	cell.SetProxyRecord(record)

	return cell, nil
}

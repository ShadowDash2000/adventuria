package cells

import (
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/helper"
	"slices"
)

type CellBase struct {
	*model.CellInfo
}

func NewCellBase(cell model.CellInfo) CellBase {
	return CellBase{&cell}
}

func (c CellBase) Data() *model.CellInfo {
	return c.CellInfo
}

func (c CellBase) Categories() []string {
	return Categories(c.Type())
}

func (c CellBase) InCategory(category string) bool {
	return slices.Contains(c.Categories(), category)
}

func (c CellBase) InCategories(categories []string) bool {
	return helper.SliceContainsAll(c.Categories(), categories)
}

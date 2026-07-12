package mocks

import (
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/helper"
	"context"
	"slices"
)

type Cell struct {
	*model.CellInfo
	CategoriesValue []string
}

func (c *Cell) Data() *model.CellInfo {
	return c.CellInfo
}

func (c *Cell) Categories() []string {
	return c.CategoriesValue
}

func (c *Cell) InCategory(category string) bool {
	return slices.Contains(c.Categories(), category)
}

func (c *Cell) InCategories(categories []string) bool {
	return helper.SliceContainsAll(c.Categories(), categories)
}

func (c *Cell) OnCellReached(context.Context, *model.Events, *model.Player, *model.ReachedContext) error {
	return nil
}

func (c *Cell) OnCellLeft(context.Context, *model.Events, *model.Player) error {
	return nil
}

package cells

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type Cell interface {
	Data() *model.CellInfo
	Categories() []string
	InCategory(category string) bool
	InCategories(categories []string) bool
}

type Rollable interface {
	Cell
	Roll(ctx context.Context, events *model.Events, player *model.Player) (*model.WheelRollResult, error)
}

type ReachedContext struct {
	Moves []*model.MoveResult
}

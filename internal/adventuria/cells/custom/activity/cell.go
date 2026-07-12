package activity

import (
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/helper"
	"context"
	"errors"
)

type activities interface {
	UpdateActivitiesFromFilter(ctx context.Context, player *model.Player, filter *model.ActivityFilter, forceUpdate bool) error
}

type filters interface {
	GetByID(ctx context.Context, id string) (*model.ActivityFilter, error)
}

var _ model.Rollable = (*CellActivity)(nil)

const Type model.CellType = "activity"

type CellActivity struct {
	cells.CellBase
	activityType model.ActivityType
	activities   activities
	filters      filters
}

func NewDef(
	activityType model.ActivityType,
	activities activities,
	activityFilters filters,
	categories ...string,
) cells.CellDef {
	return cells.NewCell(
		model.CellType(activityType),
		func(cell model.CellInfo) model.Cell {
			return &CellActivity{
				CellBase:     cells.NewCellBase(cell),
				activityType: activityType,
				activities:   activities,
				filters:      activityFilters,
			}
		},
		categories...,
	)
}

func (c *CellActivity) Roll(_ context.Context, _ *model.Events, player *model.Player) (*model.WheelRollResult, error) {
	items := player.LastAction().ItemsList()

	if len(items) == 0 {
		return nil, errors.New("no items to roll")
	}

	return &model.WheelRollResult{
		WinnerId: helper.RandomItemFromSlice(items),
	}, nil
}

func (c *CellActivity) OnCellReached(_ context.Context, _ *model.Events, _ *model.Player, _ *model.ReachedContext) error {
	return nil
}

func (c *CellActivity) OnCellLeft(_ context.Context, _ *model.Events, _ *model.Player) error {
	return nil
}

func (c *CellActivity) RefreshItems(ctx context.Context, _ *model.Events, player *model.Player) error {
	filter, err := c.filters.GetByID(ctx, c.Filter())
	if err != nil {
		return err
	}
	filter.SetType(c.activityType)
	return c.activities.UpdateActivitiesFromFilter(ctx, player, filter, true)
}

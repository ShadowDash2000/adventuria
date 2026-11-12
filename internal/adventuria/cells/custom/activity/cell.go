package activity

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/helper"
	"context"
	"errors"
)

type filters interface {
	GetByID(ctx context.Context, id string) (*model.ActivityFilterInfo, error)
}

var _ cells.Rollable = (*CellActivity)(nil)
var _ player.CellReached = (*CellActivity)(nil)

type CellActivity struct {
	cells.CellBase
	filters filters
}

func NewDef(
	activityType model.ActivityType,
	activityFilters filters,
	categories ...string,
) cells.CellDef {
	return cells.NewCell(
		model.CellType(activityType),
		func(cell model.CellInfo) (cells.Cell, error) {
			return &CellActivity{
				CellBase: cells.NewCellBase(cell),
				filters:  activityFilters,
			}, nil
		},
		categories...,
	)
}

func (c *CellActivity) Roll(_ context.Context, _ *model.Events, player *model.Player) (*model.WheelRollResult, error) {
	activitiesState := player.LastAction().State().Activities

	if activitiesState == nil {
		return nil, errs.ErrNoActiveActivity
	}

	if len(activitiesState.Ids) == 0 {
		return nil, errors.New("no items to roll")
	}

	return &model.WheelRollResult{
		WinnerId: helper.RandomItemFromSlice(activitiesState.Ids),
	}, nil
}

func (c *CellActivity) OnPlayerReached(ctx context.Context, actor *player.Actor, _ *cells.ReachedContext) error {
	filter, err := c.filters.GetByID(ctx, c.Filter())
	if err != nil {
		return err
	}

	actor.Player().LastAction().State().ActivityFilter = new(filter.Filter())

	return nil
}

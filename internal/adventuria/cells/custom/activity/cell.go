package activity

import (
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/helper"
	"context"
	"errors"
)

type activities interface {
	GetByFilter(ctx context.Context, filter model.ActivityFilter) ([]string, error)
}

type filters interface {
	GetByID(ctx context.Context, id string) (*model.ActivityFilterInfo, error)
}

var _ model.Rollable = (*CellActivity)(nil)

type CellActivity struct {
	cells.CellBase
	activities activities
	filters    filters
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
				CellBase:   cells.NewCellBase(cell),
				activities: activities,
				filters:    activityFilters,
			}
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

func (c *CellActivity) OnCellReached(ctx context.Context, _ *model.Events, player *model.Player, _ *model.ReachedContext) error {
	filter, err := c.filters.GetByID(ctx, c.Filter())
	if err != nil {
		return err
	}

	actionState := player.LastAction().State()
	actionState.ActivityFilter = &model.ActivityFilter{
		Type:            filter.Type(),
		Platforms:       filter.Platforms(),
		PlatformsStrict: filter.PlatformsStrict(),
		GameTypes:       filter.GameTypes(),
		Developers:      filter.Developers(),
		Publishers:      filter.Publishers(),
		Genres:          filter.Genres(),
		Tags:            filter.Tags(),
		Themes:          filter.Themes(),
		MinPrice:        filter.MinPrice(),
		MaxPrice:        filter.MaxPrice(),
		ReleaseDateFrom: filter.ReleaseDateFrom(),
		ReleaseDateTo:   filter.ReleaseDateTo(),
		MinCampaignTime: filter.MinCampaignTime(),
		MaxCampaignTime: filter.MaxCampaignTime(),
		Activities:      filter.Activities(),
	}
	player.LastAction().SetState(actionState)

	return nil
}

func (c *CellActivity) OnCellLeft(_ context.Context, _ *model.Events, _ *model.Player) error {
	return nil
}

func (c *CellActivity) RefreshItems(ctx context.Context, _ *model.Events, player *model.Player) error {
	actionState := player.LastAction().State()

	if actionState.ActivityFilter == nil {
		return errs.ErrNoActiveActivityFilter
	}

	ids, err := c.activities.GetByFilter(ctx, *actionState.ActivityFilter)
	if err != nil {
		return err
	}

	actionState.Activities = &model.ActionActivitiesState{
		Ids: ids,
	}
	player.LastAction().SetState(actionState)

	return nil
}

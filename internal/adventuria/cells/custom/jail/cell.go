package jail

import (
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/helper"
	"context"
	"errors"
)

type activities interface {
	GetRandomIDsByFilter(ctx context.Context, filter model.ActivityFilter) ([]string, error)
}

type filters interface {
	GetByID(ctx context.Context, id string) (*model.ActivityFilterInfo, error)
}

var _ model.Rollable = (*CellJail)(nil)

const Type model.CellType = "jail"

type CellJail struct {
	cells.CellBase
	activityType model.ActivityType
	activities   activities
	filters      filters
}

func NewDef(
	activities activities,
	activityFilters filters,
	categories ...string,
) cells.CellDef {
	return cells.NewCell(
		Type,
		func(cell model.CellInfo) model.Cell {
			return &CellJail{
				CellBase:     cells.NewCellBase(cell),
				activityType: model.ActivityTypeGame,
				activities:   activities,
				filters:      activityFilters,
			}
		},
		categories...,
	)
}

func (c *CellJail) Roll(_ context.Context, _ *model.Events, player *model.Player) (*model.WheelRollResult, error) {
	activitiesState := player.LastAction().State().Activities

	if len(activitiesState.Ids) == 0 {
		return nil, errors.New("no items to roll")
	}

	return &model.WheelRollResult{
		WinnerId: helper.RandomItemFromSlice(activitiesState.Ids),
	}, nil
}

func (c *CellJail) OnCellReached(ctx context.Context, events *model.Events, player *model.Player, _ *model.ReachedContext) error {
	if player.Progress().IsInJail() {
		player.Progress().SetCanMove(false)

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

		err = events.OnAfterGoToJail().Trigger(ctx, &model.OnAfterGoToJailEvent{})
		if err != nil {
			return err
		}
	} else {
		player.Progress().SetCanMove(true)
	}
	return nil
}

func (c *CellJail) OnCellLeft(_ context.Context, _ *model.Events, player *model.Player) error {
	// If a player somehow left a jail, we need to free them
	if player.Progress().IsInJail() {
		player.Progress().SetIsInJail(false)
		player.Progress().SetDropsInARow(0)
	}

	return nil
}

func (c *CellJail) RefreshItems(ctx context.Context, _ *model.Events, player *model.Player) error {
	actionState := player.LastAction().State()

	if actionState.ActivityFilter == nil {
		return errs.ErrNoActiveActivityFilter
	}

	ids, err := c.activities.GetRandomIDsByFilter(ctx, *actionState.ActivityFilter)
	if err != nil {
		return err
	}

	actionState.Activities = &model.ActionActivitiesState{
		Ids: ids,
	}
	player.LastAction().SetState(actionState)

	return nil
}

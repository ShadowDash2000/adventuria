package jail

import (
	"adventuria/internal/adventuria_new/cells"
	"adventuria/internal/adventuria_new/model"
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

var _ model.Rollable = (*CellJail)(nil)

const Type model.CellType = "jail"

type CellJail struct {
	cells.CellBase
	activityType model.ActivityType
	activities   activities
	filters      filters
}

func NewCellJailDef(
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

func (c *CellJail) Roll(_ context.Context, _ *model.Events, player *model.Player, _ model.RollWheelRequest) (*model.WheelRollResult, error) {
	items := player.LastAction().ItemsList()

	if len(items) == 0 {
		return nil, errors.New("no items to roll")
	}

	return &model.WheelRollResult{
		WinnerId: helper.RandomItemFromSlice(items),
	}, nil
}

func (c *CellJail) OnCellReached(_ context.Context, events *model.Events, player *model.Player, _ *model.ReachedContext) error {
	if player.Progress().IsInJail() {
		player.LastAction().SetCanMove(false)

		err := events.OnAfterGoToJail().Trigger(&model.OnAfterGoToJailEvent{})
		if err != nil {
			return err
		}
	} else {
		player.LastAction().SetCanMove(true)
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
	filter, err := c.filters.GetByID(ctx, c.Filter())
	if err != nil {
		return err
	}
	filter.SetType(c.activityType)
	return c.activities.UpdateActivitiesFromFilter(ctx, player, filter, true)
}

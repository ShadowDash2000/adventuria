package jail

import (
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/helper"
	"context"
	"errors"
)

type filters interface {
	GetByID(ctx context.Context, id string) (*model.ActivityFilterInfo, error)
}

var _ model.Rollable = (*CellJail)(nil)

const Type model.CellType = "jail"

type CellJail struct {
	cells.CellBase
	activityType model.ActivityType
	filters      filters
}

func NewDef(
	activityFilters filters,
	categories ...string,
) cells.CellDef {
	return cells.NewCell(
		Type,
		func(cell model.CellInfo) model.Cell {
			return &CellJail{
				CellBase:     cells.NewCellBase(cell),
				activityType: model.ActivityTypeGame,
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

		player.LastAction().State().ActivityFilter = new(filter.Filter())

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

package jail

import (
	"adventuria/internal/adventuria/board/actors/player"
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
var _ player.CellReached = (*CellJail)(nil)
var _ player.CellLeft = (*CellJail)(nil)

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

func (c *CellJail) OnPlayerReached(ctx context.Context, actor *player.Actor, _ *model.ReachedContext) error {
	if actor.Player().Progress().IsInJail() {
		actor.Player().Progress().SetCanMove(false)

		filter, err := c.filters.GetByID(ctx, c.Filter())
		if err != nil {
			return err
		}

		actor.Player().LastAction().State().ActivityFilter = new(filter.Filter())

		err = actor.Events().OnAfterGoToJail().Trigger(ctx, &model.OnAfterGoToJailEvent{})
		if err != nil {
			return err
		}
	} else {
		actor.Player().Progress().SetCanMove(true)
	}
	return nil
}

func (c *CellJail) OnPlayerLeft(_ context.Context, actor *player.Actor) error {
	// If a player somehow left a jail, we need to free them
	if actor.Player().Progress().IsInJail() {
		actor.Player().Progress().SetIsInJail(false)
		actor.Player().Progress().SetDropsInARow(0)
	}

	return nil
}

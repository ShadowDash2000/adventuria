package generate_wheel

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"slices"
)

type activities interface {
	GetRandomIDsByFilter(ctx context.Context, filter model.ActivityFilter) ([]string, error)
}

var _ actions.Action = (*GenerateWheel)(nil)

const Type model.ActionType = "generate_wheel"

type GenerateWheel struct {
	actions.ActionBase
	activities activities
}

func NewDef(activities activities) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() actions.Action {
			return &GenerateWheel{
				ActionBase: actions.NewActionBase(Type),
				activities: activities,
			}
		},
	)
}

func (g *GenerateWheel) CanDo(_ context.Context, actor *player.Actor) bool {
	if actor.Player().LastAction().State().ActivityFilter == nil {
		return false
	}

	return !actor.Player().Progress().CanMove() &&
		!slices.Contains([]model.ActionStatus{
			model.ActionStatusNeedToRollWheel,
			model.ActionStatusRollWheel,
		}, actor.Player().LastAction().Status())
}

func (g *GenerateWheel) Do(ctx context.Context, actor *player.Actor, _ actions.Request) (any, error) {
	actionState := actor.Player().LastAction().State()
	if actionState.ActivityFilter == nil {
		return nil, errs.ErrNoActiveActivityFilter
	}

	ids, err := g.activities.GetRandomIDsByFilter(ctx, actionState.ActivityFilter.Clone())
	if err != nil {
		return nil, err
	}

	actionState.Activities = &model.ActionActivitiesState{
		Ids: ids,
	}

	actor.Player().LastAction().SetStatus(model.ActionStatusNeedToRollWheel)

	return nil, nil
}

package generate_wheel

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"slices"
)

type activities interface {
	GetByFilter(ctx context.Context, filter model.ActivityFilter) ([]string, error)
}

var _ model.Action = (*GenerateWheel)(nil)

const Type model.ActionType = "generate_wheel"

type GenerateWheel struct {
	actions.ActionBase
	activities activities
}

func NewDef(activities activities) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &GenerateWheel{
				ActionBase: actions.NewActionBase(Type),
				activities: activities,
			}
		},
	)
}

func (g *GenerateWheel) CanDo(_ context.Context, _ *model.Events, player *model.Player) bool {
	if player.LastAction().State().ActivityFilter == nil {
		return false
	}

	return !player.Progress().CanMove() &&
		!slices.Contains([]model.ActionStatus{
			model.ActionStatusNeedToRollWheel,
			model.ActionStatusRollWheel,
		}, player.LastAction().Status())
}

func (g *GenerateWheel) Do(ctx context.Context, _ *model.Events, player *model.Player, _ model.ActionRequest) (any, error) {
	actionState := player.LastAction().State()
	if actionState.ActivityFilter == nil {
		return nil, errs.ErrNoActiveActivityFilter
	}

	ids, err := g.activities.GetByFilter(ctx, *actionState.ActivityFilter)
	if err != nil {
		return nil, err
	}

	actionState.Activities = &model.ActionActivitiesState{
		Ids: ids,
	}
	player.LastAction().SetState(actionState)
	player.LastAction().SetStatus(model.ActionStatusNeedToRollWheel)

	return nil, nil
}

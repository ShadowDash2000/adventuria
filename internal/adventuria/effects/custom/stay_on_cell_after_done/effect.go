package stay_on_cell_after_done

import (
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

type actionsService interface {
	Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
	CreateAndSetLastAction(ctx context.Context, player *model.Player, action *model.ActionInfo) (*model.ActionInfo, error)
}

type activities interface {
	GetRandomIDsByFilter(ctx context.Context, filter model.ActivityFilter) ([]string, error)
}

var _ model.Effect = (*StayOnCellAfterDone)(nil)

const Type model.EffectType = "stay_on_cell_after_done"

type StayOnCellAfterDone struct {
	effects.EffectBase
	actions    actionsService
	activities activities
}

func NewDef(actions actionsService, activities activities) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &StayOnCellAfterDone{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
				activities: activities,
			}
		},
	)
}

func (s *StayOnCellAfterDone) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (s *StayOnCellAfterDone) Subscribe(
	_ context.Context,
	events *model.Events,
	player *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnAfterDone().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterDoneEvent) error {
			lastAction := player.LastAction()
			if lastAction.Status() != model.ActionStatusDone {
				return e.Next()
			}

			actionState := player.LastAction().State()
			if actionState.ActivityFilter == nil {
				return errs.ErrNoActiveActivityFilter
			}

			lastAction, err := s.actions.Save(ctx, lastAction)
			if err != nil {
				return err
			}

			newAction, err := model.NewAction(model.ActionCreate{
				Player: player.ID(),
				Cell:   lastAction.Cell(),
				Status: model.ActionStatusRollDice,
			})
			if err != nil {
				return err
			}

			newAction, err = s.actions.CreateAndSetLastAction(ctx, player, newAction)
			if err != nil {
				return err
			}

			ids, err := s.activities.GetRandomIDsByFilter(ctx, actionState.ActivityFilter.Clone())
			if err != nil {
				return err
			}

			actionState.Activities = &model.ActionActivitiesState{
				Ids: ids,
			}

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

package stay_on_cell_after_done

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

type actionsService interface {
	Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
}

var _ effects.Effect = (*StayOnCellAfterDone)(nil)

const Type model.EffectType = "stay_on_cell_after_done"

type StayOnCellAfterDone struct {
	effects.EffectBase
	actions actionsService
}

func NewDef(actions actionsService) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &StayOnCellAfterDone{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
			}
		},
	)
}

func (s *StayOnCellAfterDone) CanUse(_ context.Context, _ *player.Actor) bool {
	return true
}

func (s *StayOnCellAfterDone) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnAfterDone().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterDoneEvent) error {
			lastAction := actor.Player().LastAction()
			if lastAction.Status() != model.ActionStatusDone {
				return e.Next()
			}

			if lastAction.State().ActivityFilter == nil {
				return errs.ErrNoActiveActivityFilter
			}

			lastAction, err := s.actions.Save(ctx, lastAction)
			if err != nil {
				return err
			}

			newAction, err := model.NewAction(model.ActionCreate{
				Player: actor.Player().ID(),
				Cell:   lastAction.Cell(),
				Status: model.ActionStatusRollDice,
			})
			if err != nil {
				return err
			}

			newActionState := lastAction.State().Clone()
			newActionState.UsedItems = nil
			newAction.SetState(newActionState)

			rootActionId := lastAction.RootAction()
			if rootActionId == "" {
				rootActionId = lastAction.ID()
			}
			newAction.SetRootAction(rootActionId)

			newAction, err = s.actions.Save(ctx, newAction)
			if err != nil {
				return err
			}

			actor.Player().SetLastAction(newAction)
			actor.Player().Progress().SetCanMove(false)

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

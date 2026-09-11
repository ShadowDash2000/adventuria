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
}

var _ model.Effect = (*StayOnCellAfterDone)(nil)

const Type model.EffectType = "stay_on_cell_after_done"

type StayOnCellAfterDone struct {
	effects.EffectBase
	actions actionsService
}

func NewDef(actions actionsService) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &StayOnCellAfterDone{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
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

			if lastAction.State().ActivityFilter == nil {
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

			player.SetLastAction(newAction)
			player.Progress().SetCanMove(false)

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

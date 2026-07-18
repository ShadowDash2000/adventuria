package stay_on_cell_after_done

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
	"errors"
)

type cellsService interface {
	GetByPlayerWrapped(ctx context.Context, player *model.Player) (model.Cell, error)
}

type actionsService interface {
	Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
}

var _ model.Effect = (*StayOnCellAfterDone)(nil)

const Type model.EffectType = "stay_on_cell_after_done"

type StayOnCellAfterDone struct {
	effects.EffectBase
	cells   cellsService
	actions actionsService
}

func NewDef(cells cellsService, actions actionsService) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &StayOnCellAfterDone{
				EffectBase: effects.NewEffectBase(effect),
				cells:      cells,
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
			if lastAction.Type() != actions.ActionTypeDone {
				return e.Next()
			}

			currentCell, err := s.cells.GetByPlayerWrapped(ctx, player)
			if err != nil {
				return err
			}

			cellRefreshable, ok := currentCell.(model.Refreshable)
			if !ok {
				return errors.New("current cell is not refreshable")
			}

			_, err = s.actions.Save(ctx, lastAction)
			if err != nil {
				return err
			}

			lastAction, err = model.NewAction(model.ActionCreate{
				Player: player.ID(),
				Cell:   lastAction.Cell(),
				Type:   actions.ActionTypeRollDice,
			})
			if err != nil {
				return err
			}

			player.SetLastAction(lastAction)

			err = cellRefreshable.RefreshItems(ctx, events, player)
			if err != nil {
				return err
			}

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

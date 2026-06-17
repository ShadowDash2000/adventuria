package safe_drop

import (
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/event_new"
	"context"
)

type cellsService interface {
	GetCurrentCellByProgress(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error)
}

var _ model.Effect = (*SafeDrop)(nil)

const Type model.EffectType = "safe_drop"

type SafeDrop struct {
	effects.EffectBase
	cells cellsService
}

func NewDef(cells cellsService) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &SafeDrop{
				EffectBase: effects.NewEffectBase(effect),
				cells:      cells,
			}
		},
	)
}

func (s *SafeDrop) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (s *SafeDrop) Subscribe(
	_ context.Context,
	events *model.Events,
	player *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event_new.Unsubscribe, error) {
	return []event_new.Unsubscribe{
		events.OnBeforeDrop().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeDropEvent) error {
			currentCell, err := s.cells.GetCurrentCellByProgress(ctx, player.Progress())
			if err != nil {
				return err
			}

			if !currentCell.Data().IsSafeDrop() {
				e.IsSafeDrop = true
				callback(ctx)
			}

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

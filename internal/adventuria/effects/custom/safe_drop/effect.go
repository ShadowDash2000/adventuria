package safe_drop

import (
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ model.Effect = (*SafeDrop)(nil)

const Type model.EffectType = "safe_drop"

type SafeDrop struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &SafeDrop{
				EffectBase: effects.NewEffectBase(effect),
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
	_ *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnActivityDrop().BindFuncWithPriority(func(ctx context.Context, e *model.OnActivityDropEvent) error {
			if !e.Result.IsSafe() {
				e.Result.SetIsSafe(true)

				if e.Mode == model.EffectRunModeApply {
					callback(ctx)
				}
			}

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

package roll_reverse

import (
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/event_new"
	"context"
)

var _ model.Effect = (*RollReverse)(nil)

const Type model.EffectType = "roll_reverse"

type RollReverse struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &RollReverse{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (r *RollReverse) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (r *RollReverse) Subscribe(
	_ context.Context,
	events *model.Events,
	_ *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event_new.Unsubscribe, error) {
	return []event_new.Unsubscribe{
		events.OnBeforeRollMove().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeRollMoveEvent) error {
			e.N *= -1
			callback(ctx)
			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

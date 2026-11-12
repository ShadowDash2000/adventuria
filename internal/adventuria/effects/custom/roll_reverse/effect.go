package roll_reverse

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ effects.Effect = (*RollReverse)(nil)

const Type model.EffectType = "roll_reverse"

type RollReverse struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &RollReverse{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (r *RollReverse) CanUse(_ context.Context, _ *player.Actor) bool {
	return true
}

func (r *RollReverse) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnBeforeRollMove().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeRollMoveEvent) error {
			e.N *= -1
			callback(ctx)
			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

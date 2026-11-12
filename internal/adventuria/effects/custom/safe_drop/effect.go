package safe_drop

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ effects.Effect = (*SafeDrop)(nil)

const Type model.EffectType = "safe_drop"

type SafeDrop struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &SafeDrop{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (s *SafeDrop) CanUse(_ context.Context, _ *player.Actor) bool {
	return true
}

func (s *SafeDrop) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnActivityDrop().BindFuncWithPriority(func(ctx context.Context, e *model.OnActivityDropEvent) error {
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

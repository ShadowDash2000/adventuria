package no_coins_for_done

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ effects.Effect = (*NoCoinsForDone)(nil)

const Type model.EffectType = "no_coins_for_done"

type NoCoinsForDone struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &NoCoinsForDone{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (n *NoCoinsForDone) CanUse(_ context.Context, _ *player.Actor) bool {
	return true
}

func (n *NoCoinsForDone) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnActivityDone().BindFuncWithPriority(func(ctx context.Context, e *model.OnActivityDoneEvent) error {
			e.Result.SetCoins(0)

			if e.Mode == model.EffectRunModeApply {
				callback(ctx)
			}

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

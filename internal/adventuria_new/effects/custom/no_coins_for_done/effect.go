package no_coins_for_done

import (
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/event_new"
	"context"
)

var _ model.Effect = (*NoCoinsForDone)(nil)

const Type model.EffectType = "no_coins_for_done"

type NoCoinsForDone struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &NoCoinsForDone{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (n *NoCoinsForDone) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (n *NoCoinsForDone) Subscribe(
	_ context.Context,
	events *model.Events,
	_ *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event_new.Unsubscribe, error) {
	return []event_new.Unsubscribe{
		events.OnBeforeDone().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeDoneEvent) error {
			e.CellCoins = 0
			callback(ctx)
			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

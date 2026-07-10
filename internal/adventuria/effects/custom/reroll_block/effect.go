package reroll_block

import (
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ model.Effect = (*RerollBlock)(nil)

const Type model.EffectType = "reroll_block"

type RerollBlock struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &RerollBlock{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (r *RerollBlock) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (r *RerollBlock) Subscribe(
	_ context.Context,
	events *model.Events,
	_ *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnBeforeRerollCheck().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeRerollCheckEvent) error {
			e.IsRerollBlocked = true
			return e.Next()
		}, effectCtx.Priority),
		events.OnAfterDone().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterDoneEvent) error {
			callback(ctx)
			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

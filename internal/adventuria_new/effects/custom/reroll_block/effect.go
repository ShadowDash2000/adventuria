package reroll_block

import (
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/event_new"
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
) ([]event_new.Unsubscribe, error) {
	return []event_new.Unsubscribe{
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

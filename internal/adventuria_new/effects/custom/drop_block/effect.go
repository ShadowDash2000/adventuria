package drop_block

import (
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/event_new"
	"context"
)

var _ model.Effect = (*DropBlock)(nil)

const Type model.EffectType = "drop_block"

type DropBlock struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &DropBlock{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (d *DropBlock) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (d *DropBlock) Subscribe(
	_ context.Context,
	events *model.Events,
	_ *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event_new.Unsubscribe, error) {
	return []event_new.Unsubscribe{
		events.OnBeforeDropCheck().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeDropCheckEvent) error {
			e.IsDropBlocked = true
			return e.Next()
		}, effectCtx.Priority),
		events.OnAfterDone().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterDoneEvent) error {
			callback(ctx)
			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

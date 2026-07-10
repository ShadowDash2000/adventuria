package debuff_block

import (
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ model.Effect = (*DebuffBlock)(nil)

const Type model.EffectType = "debuff_block"

type DebuffBlock struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &DebuffBlock{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (d *DebuffBlock) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (d *DebuffBlock) Subscribe(
	_ context.Context,
	events *model.Events,
	_ *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnBeforeItemAdd().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeItemAddEvent) error {
			if e.ItemRecord.Type() == model.ItemTypeDebuff {
				e.ShouldAddItem = false
				callback(ctx)
			}
			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

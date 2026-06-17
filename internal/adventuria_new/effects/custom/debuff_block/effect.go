package debuff_block

import (
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/event_new"
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
) ([]event_new.Unsubscribe, error) {
	return []event_new.Unsubscribe{
		events.OnBeforeItemAdd().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeItemAddEvent) error {
			if e.ItemRecord.Type() == model.ItemTypeDebuff {
				e.ShouldAddItem = false
				callback(ctx)
			}
			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

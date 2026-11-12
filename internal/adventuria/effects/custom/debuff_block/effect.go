package debuff_block

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ effects.Effect = (*DebuffBlock)(nil)

const Type model.EffectType = "debuff_block"

type DebuffBlock struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &DebuffBlock{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (d *DebuffBlock) CanUse(_ context.Context, _ *player.Actor) bool {
	return true
}

func (d *DebuffBlock) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnBeforeItemAdd().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeItemAddEvent) error {
			if e.ItemRecord.Type() == model.ItemTypeDebuff {
				e.ShouldAddItem = false
				callback(ctx)
			}
			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

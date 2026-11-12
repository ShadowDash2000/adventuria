package drop_block

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ effects.Effect = (*DropBlock)(nil)

const Type model.EffectType = "drop_block"

type DropBlock struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &DropBlock{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (d *DropBlock) CanUse(_ context.Context, _ *player.Actor) bool {
	return true
}

func (d *DropBlock) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnActivityDrop().BindFuncWithPriority(func(ctx context.Context, e *model.OnActivityDropEvent) error {
			e.Result.SetIsBlocked(true)
			return e.Next()
		}, effectCtx.Priority),
		actor.Events().OnAfterDone().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterDoneEvent) error {
			callback(ctx)
			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

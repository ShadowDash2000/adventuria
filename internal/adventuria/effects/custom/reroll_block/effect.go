package reroll_block

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ effects.Effect = (*RerollBlock)(nil)

const Type model.EffectType = "reroll_block"

type RerollBlock struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &RerollBlock{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (r *RerollBlock) CanUse(_ context.Context, _ *player.Actor) bool {
	return true
}

func (r *RerollBlock) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnBeforeRerollCheck().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeRerollCheckEvent) error {
			e.IsRerollBlocked = true
			return e.Next()
		}, effectCtx.Priority),
		actor.Events().OnAfterDone().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterDoneEvent) error {
			callback(ctx)
			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

package reset_drops_counter

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ effects.Effect = (*ResetDropCounter)(nil)

const Type model.EffectType = "reset_drops_counter"

type ResetDropCounter struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &ResetDropCounter{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (r *ResetDropCounter) CanUse(_ context.Context, actor *player.Actor) bool {
	return actor.Player().Progress().DropsInARow() > 0
}

func (r *ResetDropCounter) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnAfterItemUse(effectCtx.InvItemID).BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			actor.Player().Progress().SetDropsInARow(0)
			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

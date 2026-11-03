package reset_drops_counter

import (
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ model.Effect = (*ResetDropCounter)(nil)

const Type model.EffectType = "reset_drops_counter"

type ResetDropCounter struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &ResetDropCounter{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (r *ResetDropCounter) CanUse(_ context.Context, _ *model.Events, player *model.Player) bool {
	return player.Progress().DropsInARow() > 0
}

func (r *ResetDropCounter) Subscribe(
	_ context.Context,
	events *model.Events,
	player *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnAfterItemUse(effectCtx.InvItemID).BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			player.Progress().SetDropsInARow(0)
			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

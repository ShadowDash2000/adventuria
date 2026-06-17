package replace_dice_roll

import (
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/event_new"
	"context"
)

var _ model.Effect = (*ReplaceDiceRoll)(nil)

const Type model.EffectType = "replace_dice_roll"

type ReplaceDiceRoll struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &ReplaceDiceRoll{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (r *ReplaceDiceRoll) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (r *ReplaceDiceRoll) Subscribe(
	_ context.Context,
	events *model.Events,
	_ *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event_new.Unsubscribe, error) {
	return []event_new.Unsubscribe{
		events.OnBeforeRollMove().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeRollMoveEvent) error {
			roll, err := r.decodeValue(r.Value())
			if err != nil {
				return err
			}

			e.N = roll
			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

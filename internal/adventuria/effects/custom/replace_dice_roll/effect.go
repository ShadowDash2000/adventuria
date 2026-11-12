package replace_dice_roll

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ effects.Effect = (*ReplaceDiceRoll)(nil)

const Type model.EffectType = "replace_dice_roll"

type ReplaceDiceRoll struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &ReplaceDiceRoll{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (r *ReplaceDiceRoll) CanUse(_ context.Context, _ *player.Actor) bool {
	return true
}

func (r *ReplaceDiceRoll) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnBeforeRollMove().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeRollMoveEvent) error {
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

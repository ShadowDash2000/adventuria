package change_dices

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ effects.Effect = (*ChangeDices)(nil)

const Type model.EffectType = "change_dices"

type ChangeDices struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &ChangeDices{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (c *ChangeDices) CanUse(_ context.Context, _ *player.Actor) bool {
	return true
}

func (c *ChangeDices) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnBeforeRoll().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeRollEvent) error {
			dices, err := c.decodeValue(c.Value())
			if err != nil {
				return err
			}

			e.Dices = dices
			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

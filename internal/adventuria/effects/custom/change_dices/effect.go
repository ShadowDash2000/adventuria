package change_dices

import (
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ model.Effect = (*ChangeDices)(nil)

const Type model.EffectType = "change_dices"

type ChangeDices struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &ChangeDices{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (c *ChangeDices) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (c *ChangeDices) Subscribe(
	_ context.Context,
	events *model.Events,
	_ *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnBeforeRoll().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeRollEvent) error {
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

package drop_points_divide

import (
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ model.Effect = (*DropPointsDivide)(nil)

const Type model.EffectType = "drop_points_divide"

type DropPointsDivide struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &DropPointsDivide{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (d *DropPointsDivide) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (d *DropPointsDivide) Subscribe(
	_ context.Context,
	events *model.Events,
	_ *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnBeforeDrop().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeDropEvent) error {
			divider, err := d.decodeValue(d.Value())
			if err != nil {
				return err
			}

			e.PointsForDrop /= divider
			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

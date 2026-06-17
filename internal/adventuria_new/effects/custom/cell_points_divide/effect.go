package cell_points_divide

import (
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/event_new"
	"context"
)

var _ model.Effect = (*CellPointsDivide)(nil)

const Type model.EffectType = "cell_points_divide"

type CellPointsDivide struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &CellPointsDivide{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (c *CellPointsDivide) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (c *CellPointsDivide) Subscribe(
	_ context.Context,
	events *model.Events,
	_ *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event_new.Unsubscribe, error) {
	return []event_new.Unsubscribe{
		events.OnBeforeDone().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeDoneEvent) error {
			divider, err := c.decodeValue(c.Value())
			if err != nil {
				return err
			}

			e.CellPoints = e.CellPoints / divider
			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
		events.OnAfterMove().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterMoveEvent) error {
			callback(ctx)
			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

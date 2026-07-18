package cell_points_divide

import (
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
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
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnCompleteActivityView().BindFuncWithPriority(func(ctx context.Context, e *model.OnCompleteActivityView) error {
			divider, err := c.decodeValue(c.Value())
			if err != nil {
				return err
			}

			e.CellPoints /= divider

			return e.Next()
		}, effectCtx.Priority),
		events.OnBeforeDone().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeDoneEvent) error {
			divider, err := c.decodeValue(c.Value())
			if err != nil {
				return err
			}

			e.CellPoints /= divider
			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
		events.OnAfterMove().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterMoveEvent) error {
			callback(ctx)
			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

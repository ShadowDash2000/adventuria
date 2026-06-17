package points_change

import (
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/event_new"
	"context"
)

var _ model.Effect = (*PointsChange)(nil)

const Type model.EffectType = "points_change"

type PointsChange struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &PointsChange{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (p *PointsChange) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (p *PointsChange) Subscribe(
	_ context.Context,
	events *model.Events,
	_ *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event_new.Unsubscribe, error) {
	return []event_new.Unsubscribe{
		events.OnBeforeDone().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeDoneEvent) error {
			amount, err := p.decodeValue(p.Value())
			if err != nil {
				return err
			}

			e.CellPoints += amount
			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

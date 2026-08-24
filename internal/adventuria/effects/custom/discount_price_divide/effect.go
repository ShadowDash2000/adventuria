package discount_price_divide

import (
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ model.Effect = (*DiscountPriceDivide)(nil)

const Type model.EffectType = "discount_price_divide"

type DiscountPriceDivide struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &DiscountPriceDivide{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (d *DiscountPriceDivide) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (d *DiscountPriceDivide) Subscribe(
	_ context.Context,
	events *model.Events,
	_ *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnItemBuy().BindFuncWithPriority(func(ctx context.Context, e *model.OnItemBuyEvent) error {
			divider, err := d.decodeValue(d.Value())
			if err != nil {
				return err
			}

			e.Result.SetPrice(e.Result.Price() / divider)

			if e.Mode == model.EffectRunModeApply {
				callback(ctx)
			}

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

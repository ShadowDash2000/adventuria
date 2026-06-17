package balance_change

import (
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/event_new"
	"context"
)

var _ model.Effect = (*BalanceChange)(nil)

const Type model.EffectType = "balance_change"

type BalanceChange struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &BalanceChange{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (b *BalanceChange) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (b *BalanceChange) Subscribe(
	_ context.Context,
	events *model.Events,
	player *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event_new.Unsubscribe, error) {
	return []event_new.Unsubscribe{
		events.OnAfterItemAdd().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemAddEvent) error {
			if e.Item.Inventory().ID() != effectCtx.InvItemID {
				return e.Next()
			}

			amount, err := b.decodeValue(b.Value())
			if err != nil {
				return err
			}

			err = player.Progress().BalanceChange(amount)
			if err != nil {
				return err
			}

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

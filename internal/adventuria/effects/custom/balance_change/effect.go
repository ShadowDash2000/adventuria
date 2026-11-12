package balance_change

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ effects.Effect = (*BalanceChange)(nil)

const Type model.EffectType = "balance_change"

type BalanceChange struct {
	effects.EffectBase
}

func NewDef() effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &BalanceChange{
				EffectBase: effects.NewEffectBase(effect),
			}
		},
	)
}

func (b *BalanceChange) CanUse(_ context.Context, _ *player.Actor) bool {
	return true
}

func (b *BalanceChange) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnAfterItemAdd().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemAddEvent) error {
			if e.Item.Inventory().ID() != effectCtx.InvItemID {
				return e.Next()
			}

			amount, err := b.decodeValue(b.Value())
			if err != nil {
				return err
			}

			err = actor.Player().Progress().BalanceChange(amount)
			if err != nil {
				return err
			}

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

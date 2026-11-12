package drop_inventory

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

type inventories interface {
	DropAllItemsByPlayerID(ctx context.Context, events *model.Events, playerId string) error
}

var _ effects.Effect = (*DropInventory)(nil)

const Type model.EffectType = "drop_inventory"

type DropInventory struct {
	effects.EffectBase
	inventories inventories
}

func NewDef(inventories inventories) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &DropInventory{
				EffectBase:  effects.NewEffectBase(effect),
				inventories: inventories,
			}
		},
	)
}

func (d *DropInventory) CanUse(_ context.Context, _ *player.Actor) bool {
	return true
}

func (d *DropInventory) Subscribe(
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

			err := d.inventories.DropAllItemsByPlayerID(ctx, actor.Events(), actor.Player().ID())
			if err != nil {
				return err
			}

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

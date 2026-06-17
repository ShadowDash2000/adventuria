package drop_inventory

import (
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/event_new"
	"context"
)

type inventories interface {
	DropAllItemsByPlayerID(ctx context.Context, events *model.Events, playerId string) error
}

var _ model.Effect = (*DropInventory)(nil)

const Type model.EffectType = "drop_inventory"

type DropInventory struct {
	effects.EffectBase
	inventories inventories
}

func NewDef(inventories inventories) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &DropInventory{
				EffectBase:  effects.NewEffectBase(effect),
				inventories: inventories,
			}
		},
	)
}

func (d *DropInventory) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (d *DropInventory) Subscribe(
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

			err := d.inventories.DropAllItemsByPlayerID(ctx, events, player.ID())
			if err != nil {
				return err
			}

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

package add_items_to_inventory

import (
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/event_new"
	"context"
)

type inventories interface {
	AddItemByID(ctx context.Context, events *model.Events, player *model.Player, itemId string) (*model.InventoryItem, error)
}

type items interface {
	GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error)
}

var _ model.Effect = (*AddItemsToInventory)(nil)

const Type model.EffectType = "add_items_to_inventory"

type AddItemsToInventory struct {
	effects.EffectBase
	inventories inventories
	items       items
}

func NewDef(inventories inventories, items items) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &AddItemsToInventory{
				EffectBase:  effects.NewEffectBase(effect),
				inventories: inventories,
				items:       items,
			}
		},
	)
}

func (a *AddItemsToInventory) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (a *AddItemsToInventory) Subscribe(
	_ context.Context,
	events *model.Events,
	player *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event_new.Unsubscribe, error) {
	return []event_new.Unsubscribe{
		events.OnAfterItemUse().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			if e.InvItemId != effectCtx.InvItemID {
				return e.Next()
			}

			ids := a.decodeValue(a.Value())
			for _, id := range ids {
				_, err := a.inventories.AddItemByID(ctx, events, player, id)
				if err != nil {
					return err
				}
			}

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

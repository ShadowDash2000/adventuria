package add_items_to_inventory

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_events/custom/item_received"
	"adventuria/pkg/event"
	"context"
)

type inventories interface {
	AddItemByID(ctx context.Context, actor *player.Actor, itemId string, source item_received.Source) (*effects.InventoryItemWithEffects, error)
}

type items interface {
	GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error)
}

var _ effects.Effect = (*AddItemsToInventory)(nil)

const Type model.EffectType = "add_items_to_inventory"

type AddItemsToInventory struct {
	effects.EffectBase
	inventories inventories
	items       items
}

func NewDef(inventories inventories, items items) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &AddItemsToInventory{
				EffectBase:  effects.NewEffectBase(effect),
				inventories: inventories,
				items:       items,
			}
		},
	)
}

func (a *AddItemsToInventory) CanUse(_ context.Context, _ *player.Actor) bool {
	return true
}

func (a *AddItemsToInventory) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnAfterItemUse(effectCtx.InvItemID).BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			ids := a.decodeValue(a.Value())
			for _, id := range ids {
				_, err := a.inventories.AddItemByID(ctx, actor, id, item_received.SourceAnotherItem)
				if err != nil {
					return err
				}
			}

			callback(ctx)

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}

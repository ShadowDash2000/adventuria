package coins_for_item

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_events/custom/item_received"
	"context"
)

type inventories interface {
	AddItemByID(ctx context.Context, actor *player.Actor, itemId string, source item_received.Source) (*effects.InventoryItemWithEffects, error)
}

type items interface {
	GetByID(ctx context.Context, id string) (*model.Item, error)
}

var _ actions.Action = (*CoinsForItem)(nil)

const Type model.ActionType = "coins_for_item"

type CoinsForItem struct {
	actions.ActionBase
	inventories inventories
	items       items
}

func NewDef(inventories inventories, items items) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() actions.Action {
			return &CoinsForItem{
				ActionBase:  actions.NewActionBase(Type),
				inventories: inventories,
				items:       items,
			}
		},
	)
}

func (c *CoinsForItem) CanDo(_ context.Context, actor *player.Actor) bool {
	_, err := actor.Player().LastAction().State().Dealer.AsCoinsForItemDeal()
	if err != nil {
		return false
	}

	return true
}

func (c *CoinsForItem) Do(ctx context.Context, actor *player.Actor, _ actions.Request) (any, error) {
	deal, err := actor.Player().LastAction().State().Dealer.AsCoinsForItemDeal()
	if err != nil {
		return nil, err
	}

	_, err = c.inventories.AddItemByID(ctx, actor, deal.ItemId, item_received.SourceDealer)
	if err != nil {
		return nil, err
	}

	err = actor.Player().Progress().BalanceChange(deal.Coins)
	if err != nil {
		return nil, err
	}

	actor.Player().LastAction().State().Dealer = nil

	return nil, nil
}

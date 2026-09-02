package coins_for_item

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/model"
	"context"
)

type inventories interface {
	AddItemByID(ctx context.Context, events *model.Events, player *model.Player, itemId string) (*model.InventoryItem, error)
}

type items interface {
	GetByID(ctx context.Context, id string) (*model.Item, error)
}

var _ model.Action = (*CoinsForItem)(nil)

const Type model.ActionType = "coins_for_item"

type CoinsForItem struct {
	actions.ActionBase
	inventories inventories
	items       items
}

func NewDef(inventories inventories, items items) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &CoinsForItem{
				ActionBase:  actions.NewActionBase(Type),
				inventories: inventories,
				items:       items,
			}
		},
	)
}

func (c *CoinsForItem) CanDo(_ context.Context, _ *model.Events, player *model.Player) bool {
	_, err := player.LastAction().State().Dealer.AsCoinsForItemDeal()
	if err != nil {
		return false
	}

	return true
}

func (c *CoinsForItem) Do(ctx context.Context, events *model.Events, player *model.Player, _ model.ActionRequest) (any, error) {
	deal, err := player.LastAction().State().Dealer.AsCoinsForItemDeal()
	if err != nil {
		return nil, err
	}

	_, err = c.inventories.AddItemByID(ctx, events, player, deal.ItemId)
	if err != nil {
		return nil, err
	}

	err = player.Progress().BalanceChange(deal.Coins)
	if err != nil {
		return nil, err
	}

	player.LastAction().State().Dealer = nil

	return nil, nil
}

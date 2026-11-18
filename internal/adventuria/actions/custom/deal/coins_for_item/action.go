package coins_for_item

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_actions/custom/deals/coins_for_item"
	"adventuria/internal/adventuria/player_events/custom/item_received"
	"context"
)

type inventories interface {
	AddItemByID(ctx context.Context, actor *player.Actor, itemId string, source item_received.Source) (*effects.InventoryItemWithEffects, error)
}

type items interface {
	GetByID(ctx context.Context, id string) (*model.Item, error)
}

type playerActions interface {
	DeleteByID(ctx context.Context, id string) error
}

var _ actions.Action = (*CoinsForItem)(nil)

const Type model.ActionType = "coins_for_item"

type CoinsForItem struct {
	actions.ActionBase
	inventories   inventories
	items         items
	playerActions playerActions
}

func NewDef(inventories inventories, items items, playerActions playerActions) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() actions.Action {
			return &CoinsForItem{
				ActionBase:    actions.NewActionBase(Type),
				inventories:   inventories,
				items:         items,
				playerActions: playerActions,
			}
		},
	)
}

func (c *CoinsForItem) CanDo(_ context.Context, actor *player.Actor) bool {
	return actor.Player().Actions().Has[*coins_for_item.CoinsForItem](coins_for_item.PlayerActionType)
}

func (c *CoinsForItem) Do(ctx context.Context, actor *player.Actor, _ actions.Request) (any, error) {
	dealEntry, ok := actor.Player().Actions().Get[*coins_for_item.CoinsForItem](coins_for_item.PlayerActionType)
	if !ok {
		return nil, errs.ErrNoActiveCoinsForItemDeal
	}

	dealAction := dealEntry.Action()

	_, err := c.inventories.AddItemByID(ctx, actor, dealAction.ItemID(), item_received.SourceDealer)
	if err != nil {
		return nil, err
	}

	err = actor.Player().Progress().BalanceChange(dealAction.Coins())
	if err != nil {
		return nil, err
	}

	return nil, c.playerActions.DeleteByID(ctx, dealEntry.Data().ID())
}

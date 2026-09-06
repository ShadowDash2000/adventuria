package coins_for_item_dealer

import (
	"adventuria/internal/adventuria/action_events"
	"adventuria/internal/adventuria/model"
	"context"
)

type items interface {
	GetByID(ctx context.Context, id string) (*model.Item, error)
}

var _ model.ActionEvent = (*CoinsForItemDealer)(nil)

const Type model.ActionEventType = "coins_for_item_dealer"

type CoinsForItemDealer struct {
	action_events.ActionEventBase
	items items
}

func NewDef(items items) action_events.ActionEventDef {
	return action_events.NewActionEventDef(
		Type,
		func(cellEventInfo model.ActionEventInfo) model.ActionEvent {
			return &CoinsForItemDealer{
				ActionEventBase: action_events.NewActionEventBase(cellEventInfo),
				items:           items,
			}
		},
	)
}

func (c *CoinsForItemDealer) Init(_ context.Context, player *model.Player) error {
	actionState := player.LastAction().State()
	if actionState.Dealer != nil {
		return nil
	}

	decodedValue, err := c.decodeValue(c.Data().Value())
	if err != nil {
		return err
	}

	actionState.Dealer = model.NewCoinsForItemDeal(
		decodedValue.Coins,
		decodedValue.ItemId,
	)

	return nil
}

func (c *CoinsForItemDealer) Dispose(_ context.Context, player *model.Player) error {
	actionState := player.LastAction().State()
	if actionState.Dealer == nil {
		return nil
	}

	actionState.Dealer = nil

	return nil
}

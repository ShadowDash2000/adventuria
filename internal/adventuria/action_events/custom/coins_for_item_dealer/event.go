package coins_for_item_dealer

import (
	"adventuria/internal/adventuria/action_events"
	"adventuria/internal/adventuria/model"
	"context"
)

type items interface {
	GetByID(ctx context.Context, id string) (*model.Item, error)
}

var _ action_events.ActionEvent = (*CoinsForItemDealer)(nil)

const Type model.ActionEventType = "coins_for_item_dealer"

type CoinsForItemDealer struct {
	items   items
	payload Payload
}

func NewDef(items items) action_events.ActionEventDef {
	return action_events.NewActionEventDef(
		Type,
		func(data []byte) (action_events.ActionEvent, error) {
			return NewFromJSON(items, data)
		},
	)
}

func NewFromJSON(items items, data []byte) (*CoinsForItemDealer, error) {
	payload, err := NewPayloadFromJSON(data)
	if err != nil {
		return nil, err
	}

	return &CoinsForItemDealer{
		items:   items,
		payload: payload,
	}, nil
}

func (c *CoinsForItemDealer) Init(_ context.Context, player *model.Player) error {
	actionState := player.LastAction().State()
	if actionState.Dealer != nil {
		return nil
	}

	actionState.Dealer = model.NewCoinsForItemDeal(
		c.payload.Coins,
		c.payload.ItemId,
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

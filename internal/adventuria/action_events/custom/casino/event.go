package casino

import (
	"adventuria/internal/adventuria/action_events"
	"adventuria/internal/adventuria/model"
	"context"
)

type items interface {
	GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error)
}

var _ action_events.ActionEvent = (*Casino)(nil)

const Type model.ActionEventType = "casino"

type Casino struct {
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

func NewFromJSON(items items, data []byte) (*Casino, error) {
	payload, err := NewPayloadFromJSON(data)
	if err != nil {
		return nil, err
	}

	return &Casino{
		items:   items,
		payload: payload,
	}, nil
}

func (c *Casino) Init(_ context.Context, player *model.Player) error {
	var err error
	actionState := player.LastAction().State()
	actionState.ShopFilter = &model.ActionShopFilterState{
		Ids: c.payload.ItemIds,
	}
	actionState.Shop, err = model.NewShopState(model.ActionShopStateCreate{
		Type:            model.ShopTypeCasino,
		Ids:             c.payload.ItemIds,
		PriceMultiplier: c.payload.PriceMultiplier,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Casino) Dispose(_ context.Context, _ *model.Player) error {
	return nil
}

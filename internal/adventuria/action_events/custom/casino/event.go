package casino

import (
	"adventuria/internal/adventuria/action_events"
	"adventuria/internal/adventuria/model"
	"context"
)

type items interface {
	GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error)
}

var _ model.ActionEvent = (*Casino)(nil)

const Type model.ActionEventType = "casino"

type Casino struct {
	action_events.ActionEventBase
	items items
}

func NewDef(items items) action_events.ActionEventDef {
	return action_events.NewActionEventDef(
		Type,
		func(cellEventInfo model.ActionEventInfo) model.ActionEvent {
			return &Casino{
				ActionEventBase: action_events.NewActionEventBase(cellEventInfo),
				items:           items,
			}
		},
	)
}

func (c *Casino) Init(_ context.Context, player *model.Player) error {
	decodedValue, err := c.decodeValue(c.Data().Value())
	if err != nil {
		return err
	}

	itemsData := player.LastAction().DataList().Items
	itemsData.Ids = decodedValue.ItemIds
	itemsData.PriceMultiplier = decodedValue.PriceMultiplier
	player.LastAction().SetItemsData(itemsData)

	return nil
}

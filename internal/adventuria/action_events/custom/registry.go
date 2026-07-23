package custom

import (
	"adventuria/internal/adventuria/action_events"
	"adventuria/internal/adventuria/action_events/custom/casino"
	"adventuria/internal/adventuria/action_events/custom/coins_for_item_dealer"
	"adventuria/internal/adventuria/items"
)

func RegisterActionEvents(
	items *items.Items,
) {
	action_events.Register(
		casino.NewDef(items),
		coins_for_item_dealer.NewDef(items),
	)
}

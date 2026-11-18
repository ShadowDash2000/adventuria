package custom

import (
	"adventuria/internal/adventuria/action_events"
	"adventuria/internal/adventuria/action_events/custom/casino"
	"adventuria/internal/adventuria/action_events/custom/coins_for_item_dealer"
	"adventuria/internal/adventuria/items"
	"adventuria/internal/adventuria/player_actions"
)

func RegisterActionEvents(
	items *items.Items,
	playerActions *player_actions.PlayerActions,
) {
	action_events.Register(
		casino.NewDef(items, playerActions),
		coins_for_item_dealer.NewDef(items, playerActions),
	)
}

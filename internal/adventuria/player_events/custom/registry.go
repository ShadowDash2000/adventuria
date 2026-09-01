package custom

import (
	"adventuria/internal/adventuria/player_events"
	"adventuria/internal/adventuria/player_events/custom/item_received"
)

func RegisterPlayerEvents() {
	player_events.Register(
		item_received.NewDef(),
	)
}

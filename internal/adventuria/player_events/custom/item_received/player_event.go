package item_received

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_events"
)

var _ model.PlayerEvent = (*ItemReceived)(nil)

const Type model.PlayerEventType = "item_received"

type ItemReceived struct{}

func NewDef() player_events.PlayerEventDef {
	return player_events.NewPlayerEvent(
		Type,
		func() model.PlayerEvent {
			return &ItemReceived{}
		},
	)
}

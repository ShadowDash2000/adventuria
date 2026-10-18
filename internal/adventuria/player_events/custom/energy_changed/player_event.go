package energy_changed

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_events"
	"encoding/json"
	"fmt"
)

var _ model.PlayerEvent = (*EnergyChanged)(nil)

const Type model.PlayerEventType = "energy_changed"

type EnergyChanged struct {
	player_events.PlayerEventBase
	payload Payload
}

func NewDef() player_events.PlayerEventDef {
	return player_events.NewPlayerEvent(
		Type,
		func(payload string) (model.PlayerEvent, error) {
			return NewFromJSON(payload)
		},
	)
}

func New(payload Payload) (*EnergyChanged, error) {
	payload = payload.clone()

	err := payload.validate()
	if err != nil {
		return nil, err
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal player_event payload: %w", err)
	}

	return &EnergyChanged{
		PlayerEventBase: player_events.NewPlayerEventBase(Type, string(encoded)),
		payload:         payload,
	}, nil
}

func NewFromJSON(payload string) (*EnergyChanged, error) {
	var decoded Payload
	err := json.Unmarshal([]byte(payload), &decoded)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal player_event payload: %w", err)
	}

	err = decoded.validate()
	if err != nil {
		return nil, err
	}

	return &EnergyChanged{
		PlayerEventBase: player_events.NewPlayerEventBase(Type, payload),
		payload:         decoded,
	}, nil
}

func NewPlayerEvent(data model.PlayerEventCreate, payload Payload) (*model.PlayerEventInfo, error) {
	event, err := New(payload)
	if err != nil {
		return nil, err
	}

	return model.NewPlayerEvent(data, event)
}

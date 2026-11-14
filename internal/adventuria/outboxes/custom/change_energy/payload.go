package change_energy

import (
	"adventuria/internal/adventuria/player_events/custom/energy_changed"
	"encoding/json"
	"errors"
	"fmt"
)

type Payload struct {
	PlayerId   string                `json:"player_id"`
	ProgressId string                `json:"progress_id"`
	Amount     *int                  `json:"amount"`
	Source     energy_changed.Source `json:"source,omitempty"`
}

func NewPayloadFromJSON(data []byte) (Payload, error) {
	var payload Payload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return Payload{}, fmt.Errorf("change_energy: failed to unmarshal payload: %w", err)
	}

	err = payload.validate()
	if err != nil {
		return Payload{}, err
	}

	return payload, nil
}

func (p Payload) Marshal() ([]byte, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("change_energy: failed to marshal payload: %w", err)
	}

	return b, nil
}

func (p Payload) clone() Payload {
	var amount *int
	if p.Amount != nil {
		amount = new(*p.Amount)
	}

	return Payload{
		PlayerId:   p.PlayerId,
		ProgressId: p.ProgressId,
		Amount:     amount,
		Source:     p.Source,
	}
}

func (p Payload) validate() error {
	if p.PlayerId == "" {
		return errors.New("change_energy: player id is empty")
	}
	if p.ProgressId == "" {
		return errors.New("change_energy: progress id is empty")
	}
	if p.Amount == nil {
		return errors.New("change_energy: amount is empty")
	}
	return nil
}

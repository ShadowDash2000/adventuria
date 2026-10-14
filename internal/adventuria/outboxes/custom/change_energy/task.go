package change_energy

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/outboxes"
	"encoding/json"
	"fmt"
)

var _ model.OutboxTask = (*ChangeEnergy)(nil)

type ChangeEnergy struct {
	outboxes.OutboxTaskBase
	payload Payload
}

func New(payload Payload) (*ChangeEnergy, error) {
	payload = payload.clone()

	err := payload.validate()
	if err != nil {
		return nil, err
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("change_energy: failed to marshal payload: %w", err)
	}

	return &ChangeEnergy{
		OutboxTaskBase: outboxes.NewOutboxTaskBase(Type, string(encoded)),
		payload:        payload,
	}, nil
}

func NewFromJSON(payload string) (*ChangeEnergy, error) {
	var decoded Payload
	err := json.Unmarshal([]byte(payload), &decoded)
	if err != nil {
		return nil, fmt.Errorf("change_energy: failed to unmarshal payload: %w", err)
	}

	err = decoded.validate()
	if err != nil {
		return nil, err
	}

	return &ChangeEnergy{
		OutboxTaskBase: outboxes.NewOutboxTaskBase(Type, payload),
		payload:        decoded,
	}, nil
}

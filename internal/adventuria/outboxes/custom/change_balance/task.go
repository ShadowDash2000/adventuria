package change_balance

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/outboxes"
	"encoding/json"
	"fmt"
)

var _ model.OutboxTask = (*ChangeBalance)(nil)

type ChangeBalance struct {
	outboxes.OutboxTaskBase
	payload Payload
}

func New(payload Payload) (*ChangeBalance, error) {
	payload = payload.clone()

	err := payload.validate()
	if err != nil {
		return nil, err
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("change_balance: failed to marshal payload: %w", err)
	}

	return &ChangeBalance{
		OutboxTaskBase: outboxes.NewOutboxTaskBase(Type, string(encoded)),
		payload:        payload,
	}, nil
}

func NewFromJSON(payload string) (*ChangeBalance, error) {
	var decoded Payload
	err := json.Unmarshal([]byte(payload), &decoded)
	if err != nil {
		return nil, fmt.Errorf("change_balance: failed to unmarshal payload: %w", err)
	}

	err = decoded.validate()
	if err != nil {
		return nil, err
	}

	return &ChangeBalance{
		OutboxTaskBase: outboxes.NewOutboxTaskBase(Type, payload),
		payload:        decoded,
	}, nil
}

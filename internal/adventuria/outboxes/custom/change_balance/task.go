package change_balance

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/outboxes"
	"bytes"
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

	encoded, err := payload.Marshal()
	if err != nil {
		return nil, err
	}

	return &ChangeBalance{
		OutboxTaskBase: outboxes.NewOutboxTaskBase(Type, encoded),
		payload:        payload,
	}, nil
}

func NewFromJSON(data []byte) (*ChangeBalance, error) {
	payload, err := NewPayloadFromJSON(data)
	if err != nil {
		return nil, err
	}

	return &ChangeBalance{
		OutboxTaskBase: outboxes.NewOutboxTaskBase(Type, bytes.Clone(data)),
		payload:        payload,
	}, nil
}

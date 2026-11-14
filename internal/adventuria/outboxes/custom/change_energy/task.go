package change_energy

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/outboxes"
	"bytes"
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

	encoded, err := payload.Marshal()
	if err != nil {
		return nil, err
	}

	return &ChangeEnergy{
		OutboxTaskBase: outboxes.NewOutboxTaskBase(Type, encoded),
		payload:        payload,
	}, nil
}

func NewFromJSON(data []byte) (*ChangeEnergy, error) {
	payload, err := NewPayloadFromJSON(data)
	if err != nil {
		return nil, err
	}

	return &ChangeEnergy{
		OutboxTaskBase: outboxes.NewOutboxTaskBase(Type, bytes.Clone(data)),
		payload:        payload,
	}, nil
}

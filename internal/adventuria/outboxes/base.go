package outboxes

import (
	"adventuria/internal/adventuria/model"
	"bytes"
)

type OutboxTaskBase struct {
	t           model.OutboxType
	payloadJSON []byte
}

func NewOutboxTaskBase(t model.OutboxType, data []byte) OutboxTaskBase {
	return OutboxTaskBase{
		t:           t,
		payloadJSON: data,
	}
}

func (o OutboxTaskBase) Type() model.OutboxType {
	return o.t
}

func (o OutboxTaskBase) Payload() []byte {
	return bytes.Clone(o.payloadJSON)
}

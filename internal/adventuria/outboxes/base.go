package outboxes

import "adventuria/internal/adventuria/model"

type OutboxTaskBase struct {
	t           model.OutboxType
	payloadJSON string
}

func NewOutboxTaskBase(t model.OutboxType, payload string) OutboxTaskBase {
	return OutboxTaskBase{
		t:           t,
		payloadJSON: payload,
	}
}

func (o OutboxTaskBase) Type() model.OutboxType {
	return o.t
}

func (o OutboxTaskBase) Payload() string {
	return o.payloadJSON
}

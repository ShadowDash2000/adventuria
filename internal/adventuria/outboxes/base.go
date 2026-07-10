package outboxes

import "adventuria/internal/adventuria/model"

type OutboxBase struct {
	t model.OutboxType
}

func NewOutboxBase(t model.OutboxType) OutboxBase {
	return OutboxBase{t: t}
}

func (a OutboxBase) Type() model.OutboxType {
	return a.t
}

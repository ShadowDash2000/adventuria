package outboxes

import (
	"adventuria/internal/adventuria/model"
)

type OutboxDef struct {
	t   model.OutboxType
	new OutboxCreator
}

func (o *OutboxDef) New() model.Outbox {
	return o.new()
}

func (o *OutboxDef) Type() model.OutboxType {
	return o.t
}

type OutboxCreator func() model.Outbox

var registry = &Registry{outboxes: map[model.OutboxType]OutboxDef{}}

type Registry struct {
	outboxes map[model.OutboxType]OutboxDef
}

func (r *Registry) Register(outboxes ...OutboxDef) {
	for _, outbox := range outboxes {
		r.outboxes[outbox.t] = outbox
	}
}

func (r *Registry) Get(t model.OutboxType) (OutboxDef, bool) {
	e, ok := r.outboxes[t]
	return e, ok
}

func NewOutbox(t model.OutboxType, new OutboxCreator) OutboxDef {
	return OutboxDef{
		t:   t,
		new: new,
	}
}

func Register(outboxes ...OutboxDef) {
	registry.Register(outboxes...)
}

func Get(t model.OutboxType) (OutboxDef, bool) {
	return registry.Get(t)
}

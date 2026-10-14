package outboxes

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type OutboxDef struct {
	t       model.OutboxType
	process ProcessFunc
}

func (o *OutboxDef) Type() model.OutboxType {
	return o.t
}

func (o *OutboxDef) Process(ctx context.Context, payload string) error {
	return o.process(ctx, payload)
}

type ProcessFunc func(ctx context.Context, payload string) error

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

func NewOutbox(t model.OutboxType, process ProcessFunc) OutboxDef {
	return OutboxDef{
		t:       t,
		process: process,
	}
}

func Register(outboxes ...OutboxDef) {
	registry.Register(outboxes...)
}

func Get(t model.OutboxType) (OutboxDef, bool) {
	return registry.Get(t)
}

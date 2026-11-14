package action_events

import (
	"adventuria/internal/adventuria/model"
)

type ActionEventDef struct {
	t model.ActionEventType
	c ActionEventCreator
}

func (a *ActionEventDef) New(actionEvent *model.ActionEventInfo) (ActionEvent, error) {
	return a.c([]byte(actionEvent.Value()))
}

type ActionEventCreator func(data []byte) (ActionEvent, error)

var registry = &Registry{actionEvents: map[model.ActionEventType]ActionEventDef{}}

type Registry struct {
	actionEvents map[model.ActionEventType]ActionEventDef
}

func (r *Registry) Register(actionEvents ...ActionEventDef) {
	for _, actionEvent := range actionEvents {
		r.actionEvents[actionEvent.t] = actionEvent
	}
}

func (r *Registry) Get(t model.ActionEventType) (ActionEventDef, bool) {
	e, ok := r.actionEvents[t]
	return e, ok
}

func NewActionEventDef(t model.ActionEventType, c ActionEventCreator) ActionEventDef {
	return ActionEventDef{
		t: t,
		c: c,
	}
}

func Register(actionEventDefs ...ActionEventDef) {
	registry.Register(actionEventDefs...)
}

func Get(t model.ActionEventType) (ActionEventDef, bool) {
	return registry.Get(t)
}

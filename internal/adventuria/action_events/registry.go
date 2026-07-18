package action_events

import (
	"adventuria/internal/adventuria/model"
	"fmt"
)

type ActionEventDef struct {
	t   model.ActionEventType
	new ActionEventCreator
}

type ActionEventCreator func(actionEventInfo model.ActionEventInfo) model.ActionEvent

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

func NewActionEventDef(t model.ActionEventType, new ActionEventCreator) ActionEventDef {
	return ActionEventDef{
		t:   t,
		new: new,
	}
}

func Register(actionEventDefs ...ActionEventDef) {
	registry.Register(actionEventDefs...)
}

func Get(t model.ActionEventType) (ActionEventDef, bool) {
	return registry.Get(t)
}

func Create(actionEventInfo model.ActionEventInfo) (model.ActionEvent, error) {
	actionEventDef, ok := Get(actionEventInfo.Type())
	if !ok {
		return nil, fmt.Errorf("action event type %s not registered", actionEventInfo.Type())
	}
	return actionEventDef.new(actionEventInfo), nil
}

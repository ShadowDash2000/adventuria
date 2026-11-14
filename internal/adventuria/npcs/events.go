package npcs

import "adventuria/pkg/event"

type UpdateStartEvent struct {
	event.Event
}

func (n *NPCs) OnUpdateStart() *event.Hook[*UpdateStartEvent] {
	return n.onUpdateStart
}

type UpdateEndEvent struct {
	event.Event
}

func (n *NPCs) OnUpdateEnd() *event.Hook[*UpdateEndEvent] {
	return n.onUpdateEnd
}

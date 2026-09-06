package cell_events_schedules

import "adventuria/pkg/event"

type UpdateStartEvent struct {
	event.Event
}

func (c *CellEventsSchedules) OnUpdateStart() *event.Hook[*UpdateStartEvent] {
	return c.onUpdateStart
}

type UpdateEndEvent struct {
	event.Event
}

func (c *CellEventsSchedules) OnUpdateEnd() *event.Hook[*UpdateEndEvent] {
	return c.onUpdateEnd
}

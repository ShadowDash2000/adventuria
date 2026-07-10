package adventuria

import "adventuria/pkg/event"

type onKillParserEvent struct {
	event.Event
}

func (g *Game) onKillParser() *event.Hook[*onKillParserEvent] {
	return g.onKillParserEvent
}

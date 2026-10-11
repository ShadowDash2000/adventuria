package player_events

import "adventuria/internal/adventuria/model"

type PlayerEventDef struct {
	t           model.PlayerEventType
	newFromJSON PlayerEventCreator
}

func (p *PlayerEventDef) NewFromJSON(payload string) (model.PlayerEvent, error) {
	return p.newFromJSON(payload)
}

func (p *PlayerEventDef) Type() model.PlayerEventType {
	return p.t
}

type PlayerEventCreator func(string) (model.PlayerEvent, error)

var registry = &Registry{playerEvents: map[model.PlayerEventType]PlayerEventDef{}}

type Registry struct {
	playerEvents map[model.PlayerEventType]PlayerEventDef
}

func (r *Registry) Register(playerEvents ...PlayerEventDef) {
	for _, playerEvent := range playerEvents {
		r.playerEvents[playerEvent.t] = playerEvent
	}
}

func (r *Registry) Get(t model.PlayerEventType) (PlayerEventDef, bool) {
	e, ok := r.playerEvents[t]
	return e, ok
}

func NewPlayerEvent(t model.PlayerEventType, newFromJSON PlayerEventCreator) PlayerEventDef {
	return PlayerEventDef{
		t:           t,
		newFromJSON: newFromJSON,
	}
}

func Register(playerEvents ...PlayerEventDef) {
	registry.Register(playerEvents...)
}

func Get(t model.PlayerEventType) (PlayerEventDef, bool) {
	return registry.Get(t)
}

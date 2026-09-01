package player_events

import "adventuria/internal/adventuria/model"

type PlayerEventDef struct {
	t   model.PlayerEventType
	new PlayerEventCreator
}

func (p *PlayerEventDef) New() model.PlayerEvent {
	return p.new()
}

func (p *PlayerEventDef) Type() model.PlayerEventType {
	return p.t
}

type PlayerEventCreator func() model.PlayerEvent

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

func NewPlayerEvent(t model.PlayerEventType, new PlayerEventCreator) PlayerEventDef {
	return PlayerEventDef{
		t:   t,
		new: new,
	}
}

func Register(playerEvents ...PlayerEventDef) {
	registry.Register(playerEvents...)
}

func Get(t model.PlayerEventType) (PlayerEventDef, bool) {
	return registry.Get(t)
}

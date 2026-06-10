package scope

import (
	"adventuria/internal/adventuria_new/model"
)

type Scope struct {
	events *model.Events
	player *model.Player
}

func New(
	player *model.Player,
) *Scope {
	return &Scope{
		events: model.NewEvents(),
		player: player,
	}
}

func (e *Scope) Events() *model.Events {
	return e.events
}

func (e *Scope) Player() *model.Player {
	return e.player
}

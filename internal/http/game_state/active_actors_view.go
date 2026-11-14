package game_state

import (
	"adventuria/internal/adventuria/board/actors"
	"time"
)

type activeActorView struct {
	Id           string           `json:"id"`
	Kind         actors.ActorKind `json:"kind"`
	Updated      time.Time        `json:"updated"`
	Name         string           `json:"name"`
	Avatar       string           `json:"avatar"`
	Color        string           `json:"color"`
	CurrentWorld string           `json:"current_world"`
	CellsPassed  int              `json:"cells_passed"`
}

func activeActorToView(actor *actors.ActiveActor) *activeActorView {
	return &activeActorView{
		Id:           actor.Id,
		Kind:         actor.Kind,
		Updated:      actor.Updated,
		Name:         actor.Name,
		Avatar:       actor.Avatar,
		Color:        actor.Color,
		CurrentWorld: actor.CurrentWorld,
		CellsPassed:  actor.CellsPassed,
	}
}

func activeActorsToView(activeActors []*actors.ActiveActor) []*activeActorView {
	res := make([]*activeActorView, len(activeActors))
	for i, actor := range activeActors {
		res[i] = activeActorToView(actor)
	}
	return res
}

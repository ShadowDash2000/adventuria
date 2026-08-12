package players

import "adventuria/internal/adventuria/model"

type CompletedActivity struct {
	Player CompletedActivityPlayer
	Status model.ActionStatus
}

type CompletedActivityPlayer struct {
	Id     string
	Name   string
	Avatar string
	Color  string
}

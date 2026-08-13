package completed_activities

import "adventuria/internal/adventuria/model"

type CompletedActivity struct {
	Id       string
	Name     string
	Cover    string
	CoverAlt string
	Players  []*CompletedActivityPlayerStatus
}

type CompletedActivityPlayerStatus struct {
	Player *CompletedActivityPlayer
	Status model.ActionStatus
}

type CompletedActivityPlayer struct {
	Id     string
	Name   string
	Avatar string
	Color  string
}

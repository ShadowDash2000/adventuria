package completed_activities

import (
	"adventuria/internal/adventuria/model"
	"time"
)

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
	Date   time.Time
}

type CompletedActivityPlayer struct {
	Id     string
	Name   string
	Avatar string
	Color  string
}

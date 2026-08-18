package completed_activities

import (
	"adventuria/internal/adventuria/model"
	"time"
)

type completedActivityView struct {
	CollectionName string              `json:"collectionName"`
	Id             string              `json:"id"`
	Name           string              `json:"name"`
	Cover          string              `json:"cover"`
	CoverAlt       string              `json:"cover_alt"`
	Players        []*playerStatusView `json:"players"`
}

type playerStatusView struct {
	Player *playerView        `json:"player"`
	Status model.ActionStatus `json:"status"`
	Date   time.Time          `json:"date"`
}

type playerView struct {
	CollectionName string `json:"collectionName"`
	Id             string `json:"id"`
	Name           string `json:"name"`
	Avatar         string `json:"avatar"`
	Color          string `json:"color"`
}

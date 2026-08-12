package completed_activities

import "adventuria/internal/adventuria/model"

type completedActivity struct {
	Player player             `json:"player"`
	Status model.ActionStatus `json:"status"`
}

type player struct {
	CollectionName string `json:"collectionName"`
	Id             string `json:"id"`
	Name           string `json:"name"`
	Avatar         string `json:"avatar"`
	Color          string `json:"color"`
}

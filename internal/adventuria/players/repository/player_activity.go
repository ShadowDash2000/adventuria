package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/players"
)

type playerActivityRow struct {
	ActivityId string `db:"activity_id"`
	Status     string `db:"status"`
	PlayerId   string `db:"player_id"`
	Name       string `db:"name"`
	Avatar     string `db:"avatar"`
	Color      string `db:"color"`
}

func playerActivityRowToCompletedActivity(playerActivity *playerActivityRow) *players.CompletedActivity {
	return &players.CompletedActivity{
		Player: players.CompletedActivityPlayer{
			Id:     playerActivity.PlayerId,
			Name:   playerActivity.Name,
			Avatar: playerActivity.Avatar,
			Color:  playerActivity.Color,
		},
		Status: model.ActionStatus(playerActivity.Status),
	}
}

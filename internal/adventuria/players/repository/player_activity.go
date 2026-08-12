package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/players"
)

type playerActivityRow struct {
	ActivityId       string `db:"activity_id"`
	ActivityName     string `db:"activity_name"`
	ActivityCover    string `db:"activity_cover"`
	ActivityCoverAlt string `db:"activity_cover_alt"`
	Status           string `db:"status"`
	PlayerId         string `db:"player_id"`
	PlayerName       string `db:"player_name"`
	PlayerAvatar     string `db:"player_avatar"`
	PlayerColor      string `db:"player_color"`
}

func playerActivityRowToCompletedActivity(playerActivity *playerActivityRow) *players.CompletedActivity {
	return &players.CompletedActivity{
		Id:       playerActivity.ActivityId,
		Name:     playerActivity.ActivityName,
		Cover:    playerActivity.ActivityCover,
		CoverAlt: playerActivity.ActivityCoverAlt,
		Players:  nil,
	}
}

func playerActivityRowToCompletedActivityPlayerStatus(playerActivity *playerActivityRow) *players.CompletedActivityPlayerStatus {
	return &players.CompletedActivityPlayerStatus{
		Player: playerActivityRowToCompletedActivityPlayer(playerActivity),
		Status: model.ActionStatus(playerActivity.Status),
	}
}

func playerActivityRowToCompletedActivityPlayer(playerActivity *playerActivityRow) *players.CompletedActivityPlayer {
	return &players.CompletedActivityPlayer{
		Id:     playerActivity.PlayerId,
		Name:   playerActivity.PlayerName,
		Avatar: playerActivity.PlayerAvatar,
		Color:  playerActivity.PlayerColor,
	}
}

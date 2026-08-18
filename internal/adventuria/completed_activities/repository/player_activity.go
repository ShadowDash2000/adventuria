package repository

import (
	"adventuria/internal/adventuria/completed_activities"
	"adventuria/internal/adventuria/model"

	"github.com/pocketbase/pocketbase/tools/types"
)

type playerActivityRow struct {
	ActivityId       string         `db:"activity_id"`
	ActivityName     string         `db:"activity_name"`
	ActivityCover    string         `db:"activity_cover"`
	ActivityCoverAlt string         `db:"activity_cover_alt"`
	ActionCreated    types.DateTime `db:"created"`
	Status           string         `db:"status"`
	PlayerId         string         `db:"player_id"`
	PlayerName       string         `db:"player_name"`
	PlayerAvatar     string         `db:"player_avatar"`
	PlayerColor      string         `db:"player_color"`
}

func playerActivityRowToCompletedActivity(playerActivity *playerActivityRow) *completed_activities.CompletedActivity {
	return &completed_activities.CompletedActivity{
		Id:       playerActivity.ActivityId,
		Name:     playerActivity.ActivityName,
		Cover:    playerActivity.ActivityCover,
		CoverAlt: playerActivity.ActivityCoverAlt,
		Players:  nil,
	}
}

func playerActivityRowToCompletedActivityPlayerStatus(playerActivity *playerActivityRow) *completed_activities.CompletedActivityPlayerStatus {
	return &completed_activities.CompletedActivityPlayerStatus{
		Player: playerActivityRowToCompletedActivityPlayer(playerActivity),
		Status: model.ActionStatus(playerActivity.Status),
		Date:   playerActivity.ActionCreated.Time(),
	}
}

func playerActivityRowToCompletedActivityPlayer(playerActivity *playerActivityRow) *completed_activities.CompletedActivityPlayer {
	return &completed_activities.CompletedActivityPlayer{
		Id:     playerActivity.PlayerId,
		Name:   playerActivity.PlayerName,
		Avatar: playerActivity.PlayerAvatar,
		Color:  playerActivity.PlayerColor,
	}
}

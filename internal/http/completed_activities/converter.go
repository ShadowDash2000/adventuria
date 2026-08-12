package completed_activities

import (
	"adventuria/internal/adventuria/players"
	"adventuria/internal/adventuria/schema"
)

func completedActivityToView(activity *players.CompletedActivity) *completedActivity {
	return &completedActivity{
		Player: player{
			CollectionName: schema.CollectionPlayers,
			Id:             activity.Player.Id,
			Name:           activity.Player.Name,
			Avatar:         activity.Player.Avatar,
			Color:          activity.Player.Color,
		},
		Status: activity.Status,
	}
}

func completedActivityMapToViewMap(activities map[string][]*players.CompletedActivity) map[string][]*completedActivity {
	res := make(map[string][]*completedActivity, len(activities))
	for activityId, completedActivities := range activities {
		res[activityId] = make([]*completedActivity, len(completedActivities))
		for i, activity := range completedActivities {
			res[activityId][i] = completedActivityToView(activity)
		}
	}
	return res
}

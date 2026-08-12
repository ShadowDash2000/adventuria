package completed_activities

import (
	"adventuria/internal/adventuria/players"
	"adventuria/internal/adventuria/schema"
)

func completedActivityToView(activity *players.CompletedActivity) *completedActivityView {
	playersView := make([]*playerStatusView, len(activity.Players))
	for i, activityPlayer := range activity.Players {
		playersView[i] = completedActivityPlayerStatusToView(activityPlayer)
	}

	return &completedActivityView{
		CollectionName: schema.CollectionActivities,
		Id:             activity.Id,
		Name:           activity.Name,
		Cover:          activity.Cover,
		CoverAlt:       activity.CoverAlt,
		Players:        playersView,
	}
}

func completedActivitiesToView(completedActivities []*players.CompletedActivity) []*completedActivityView {
	res := make([]*completedActivityView, len(completedActivities))
	for i, activity := range completedActivities {
		res[i] = completedActivityToView(activity)
	}
	return res
}

func completedActivityPlayerStatusToView(activityPlayerStatus *players.CompletedActivityPlayerStatus) *playerStatusView {
	return &playerStatusView{
		Player: completedActivityPlayerToView(activityPlayerStatus.Player),
		Status: activityPlayerStatus.Status,
	}
}

func completedActivityPlayerToView(activityPlayer *players.CompletedActivityPlayer) *playerView {
	return &playerView{
		CollectionName: schema.CollectionPlayers,
		Id:             activityPlayer.Id,
		Name:           activityPlayer.Name,
		Avatar:         activityPlayer.Avatar,
		Color:          activityPlayer.Color,
	}
}

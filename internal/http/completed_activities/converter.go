package completed_activities

import (
	"adventuria/internal/adventuria/completed_activities"
)

func completedActivityToView(activity *completed_activities.CompletedActivity) *completedActivityView {
	playersView := make([]*playerStatusView, len(activity.Players))
	for i, activityPlayer := range activity.Players {
		playersView[i] = completedActivityPlayerStatusToView(activityPlayer)
	}

	return &completedActivityView{
		Id:       activity.Id,
		Name:     activity.Name,
		Cover:    activity.Cover,
		CoverAlt: activity.CoverAlt,
		Players:  playersView,
	}
}

func completedActivitiesToView(completedActivities []*completed_activities.CompletedActivity) []*completedActivityView {
	res := make([]*completedActivityView, len(completedActivities))
	for i, activity := range completedActivities {
		res[i] = completedActivityToView(activity)
	}
	return res
}

func completedActivityPlayerStatusToView(activityPlayerStatus *completed_activities.CompletedActivityPlayerStatus) *playerStatusView {
	return &playerStatusView{
		Player: completedActivityPlayerToView(activityPlayerStatus.Player),
		Status: activityPlayerStatus.Status,
		Date:   activityPlayerStatus.Date,
	}
}

func completedActivityPlayerToView(activityPlayer *completed_activities.CompletedActivityPlayer) *playerView {
	return &playerView{
		Id:     activityPlayer.Id,
		Name:   activityPlayer.Name,
		Avatar: activityPlayer.Avatar,
		Color:  activityPlayer.Color,
	}
}

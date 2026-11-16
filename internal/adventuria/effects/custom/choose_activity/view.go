package choose_activity

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
)

var _ effects.View = (*ChooseActivity)(nil)

type activityView struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (c *ChooseActivity) GetView(ctx context.Context, actor *player.Actor) (any, error) {
	activitiesState := actor.Player().LastAction().State().Activities
	if activitiesState == nil {
		return nil, errs.ErrNoActiveActivity
	}

	activities, err := c.activities.GetByIDs(ctx, activitiesState.Ids)
	if err != nil {
		return nil, err
	}

	return activitiesToActivitiesView(activities), nil
}

func activityToActivityView(activity *model.Activity) *activityView {
	return &activityView{
		Id:   activity.ID(),
		Name: activity.Name(),
	}
}

func activitiesToActivitiesView(activities []*model.Activity) []*activityView {
	activitiesView := make([]*activityView, len(activities))
	for i, activity := range activities {
		activitiesView[i] = activityToActivityView(activity)
	}
	return activitiesView
}

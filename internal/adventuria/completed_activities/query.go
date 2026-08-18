package completed_activities

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type repository interface {
	GetActivitiesByCellIDWithStatus(ctx context.Context, cellId string, statuses []model.ActionStatus, limitActivities, limitActions int) ([]*CompletedActivity, error)
}

type Query struct {
	repository repository
}

func NewQuery(repository repository) *Query {
	return &Query{
		repository: repository,
	}
}

func (q *Query) GetCompletedActivitiesByCellID(ctx context.Context, cellId string) ([]*CompletedActivity, error) {
	return q.repository.GetActivitiesByCellIDWithStatus(
		ctx,
		cellId,
		[]model.ActionStatus{
			model.ActionStatusDone,
			model.ActionStatusDrop,
		},
		20,
		100,
	)
}

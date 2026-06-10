package activity_filters

import (
	"adventuria/internal/adventuria_new/model"
	"context"
)

type repository interface {
	GetByID(ctx context.Context, id string) (*model.ActivityFilter, error)
}

type ActivityFilters struct {
	repository repository
}

func NewActivityFilters(repository repository) *ActivityFilters {
	return &ActivityFilters{repository: repository}
}

func (a *ActivityFilters) GetByID(ctx context.Context, id string) (*model.ActivityFilter, error) {
	return a.repository.GetByID(ctx, id)
}

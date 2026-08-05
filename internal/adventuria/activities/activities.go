package activities

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type repository interface {
	Create(ctx context.Context, activity *model.Activity) (*model.Activity, error)
	Update(ctx context.Context, activity *model.Activity) (*model.Activity, error)
	GetByIdDb(ctx context.Context, idDb string) (*model.Activity, error)
	GetByFilter(ctx context.Context, filter model.ActivityFilter, poolSize, resultSize int) ([]string, error)
	GetByID(ctx context.Context, id string) (*model.Activity, error)
	GetByIDs(ctx context.Context, ids []string) ([]*model.Activity, error)
	GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error)
}

type Activities struct {
	repository repository
}

func NewActivities(repository repository) *Activities {
	return &Activities{repository: repository}
}

func (a *Activities) GetOrCreate(ctx context.Context, data model.ActivityCreate) (*model.Activity, error) {
	activity, err := a.repository.GetByIdDb(ctx, data.IdDb)
	if err != nil {
		if errors.Is(err, errs.ErrActivityNotFound) {
			return model.NewActivity(data)
		}
		return nil, err
	}

	return activity, nil
}

func (a *Activities) GetByFilter(ctx context.Context, filter model.ActivityFilter) ([]string, error) {
	return a.repository.GetByFilter(ctx, filter, 20000, 20)
}

func (a *Activities) GetByID(ctx context.Context, id string) (*model.Activity, error) {
	return a.repository.GetByID(ctx, id)
}

func (a *Activities) GetByIDs(ctx context.Context, ids []string) ([]*model.Activity, error) {
	return a.repository.GetByIDs(ctx, ids)
}

func (a *Activities) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	return a.repository.GetChecksumsByIDs(ctx, ids)
}

func (a *Activities) Save(ctx context.Context, activity *model.Activity) (*model.Activity, error) {
	if activity.IsNew() {
		return a.repository.Create(ctx, activity)
	}

	return a.repository.Update(ctx, activity)
}

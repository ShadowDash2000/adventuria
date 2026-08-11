package activities

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
	"math/rand/v2"
)

type repository interface {
	Create(ctx context.Context, activity *model.Activity) (*model.Activity, error)
	Update(ctx context.Context, activity *model.Activity) (*model.Activity, error)
	GetByIdDb(ctx context.Context, idDb string) (*model.Activity, error)
	GetCountByFilter(ctx context.Context, filter model.ActivityFilter) (int, error)
	GetIDsByFilter(ctx context.Context, filter model.ActivityFilter, offset, limit int) ([]string, error)
	GetAverageCampaignTimeByFilter(ctx context.Context, filter model.ActivityFilter) (float64, error)
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

func (a *Activities) Save(ctx context.Context, activity *model.Activity) (*model.Activity, error) {
	if activity.IsNew() {
		return a.repository.Create(ctx, activity)
	}

	return a.repository.Update(ctx, activity)
}

const (
	poolSize   = 20000
	resultSize = 20
)

func (a *Activities) GetRandomIDsByFilter(ctx context.Context, filter model.ActivityFilter) ([]string, error) {
	count, err := a.repository.GetCountByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	limit := count
	offset := 0

	if count > poolSize {
		limit = poolSize
		offset = rand.N(count - poolSize + 1)
	}

	ids, err := a.GetIDsByFilter(ctx, filter, offset, limit)
	if err != nil {
		return nil, err
	}

	rand.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})

	resultSize := min(len(ids), resultSize)

	res := make([]string, resultSize)
	for i := range resultSize {
		res[i] = ids[i]
	}

	return res, nil
}

func (a *Activities) GetIDsByFilter(ctx context.Context, filter model.ActivityFilter, offset, limit int) ([]string, error) {
	return a.repository.GetIDsByFilter(ctx, filter, offset, limit)
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

func (a *Activities) GetAverageCampaignTimeByFilter(ctx context.Context, filter model.ActivityFilter) (float64, error) {
	return a.repository.GetAverageCampaignTimeByFilter(ctx, filter)
}

package mocks

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type Activities struct {
	GetByFilterFunc func(ctx context.Context, filter model.ActivityFilter) ([]string, error)
	GetByIDFunc     func(ctx context.Context, id string) (*model.Activity, error)
	GetByIDsFunc    func(ctx context.Context, ids []string) ([]*model.Activity, error)
}

func (m *Activities) GetByFilter(ctx context.Context, filter model.ActivityFilter) ([]string, error) {
	if m.GetByFilterFunc != nil {
		return m.GetByFilterFunc(ctx, filter)
	}

	return nil, nil
}

func (m *Activities) GetByID(ctx context.Context, id string) (*model.Activity, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}

	return nil, nil
}

func (m *Activities) GetByIDs(ctx context.Context, ids []string) ([]*model.Activity, error) {
	if m.GetByIDsFunc != nil {
		return m.GetByIDsFunc(ctx, ids)
	}

	return nil, nil
}

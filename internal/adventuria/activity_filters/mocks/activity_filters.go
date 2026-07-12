package mocks

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type ActivityFilters struct {
	GetByIDFunc func(ctx context.Context, id string) (*model.ActivityFilter, error)
}

func (m *ActivityFilters) GetByID(ctx context.Context, id string) (*model.ActivityFilter, error) {
	if m.GetByIDFunc == nil {
		return nil, nil
	}

	return m.GetByIDFunc(ctx, id)
}

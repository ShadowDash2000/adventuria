package mocks

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type Items struct {
	GetByIDsFunc func(ctx context.Context, ids []string) ([]*model.Item, error)
}

func (m *Items) GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error) {
	if m.GetByIDsFunc == nil {
		return nil, nil
	}

	return m.GetByIDsFunc(ctx, ids)
}

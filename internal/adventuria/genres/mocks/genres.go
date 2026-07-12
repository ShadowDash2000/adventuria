package mocks

import (
	"context"
)

type Genres struct {
	ExistsFunc func(ctx context.Context, id string) (bool, error)
}

func (m *Genres) Exists(ctx context.Context, id string) (bool, error) {
	if m.ExistsFunc == nil {
		return false, nil
	}

	return m.ExistsFunc(ctx, id)
}

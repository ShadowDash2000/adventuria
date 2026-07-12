package mocks

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type Cells struct {
	GetCurrentCellByProgressFunc func(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error)
}

func (m *Cells) GetCurrentCellByProgress(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error) {
	if m.GetCurrentCellByProgressFunc == nil {
		return nil, nil
	}

	return m.GetCurrentCellByProgressFunc(ctx, progress)
}

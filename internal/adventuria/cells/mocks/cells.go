package mocks

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type Cells struct {
	GetByPlayerFunc        func(ctx context.Context, player *model.Player) (*model.CellInfo, error)
	GetByPlayerWrappedFunc func(ctx context.Context, player *model.Player) (model.Cell, error)
}

func (m *Cells) GetByPlayer(ctx context.Context, player *model.Player) (*model.CellInfo, error) {
	if m.GetByPlayerFunc == nil {
		return nil, nil
	}

	return m.GetByPlayerFunc(ctx, player)
}

func (m *Cells) GetByPlayerWrapped(ctx context.Context, player *model.Player) (model.Cell, error) {
	if m.GetByPlayerWrappedFunc == nil {
		return nil, nil
	}

	return m.GetByPlayerWrappedFunc(ctx, player)
}

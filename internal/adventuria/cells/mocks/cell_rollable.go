package mocks

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type RollableCell struct {
	*Cell
	RefreshCalled bool
}

func (m *RollableCell) RefreshItems(_ context.Context, _ *model.Events, _ *model.Player) error {
	m.RefreshCalled = true
	return nil
}

func (m *RollableCell) Roll(_ context.Context, _ *model.Events, _ *model.Player) (*model.WheelRollResult, error) {
	return nil, nil
}

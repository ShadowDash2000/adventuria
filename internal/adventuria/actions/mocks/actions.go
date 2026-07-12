package mocks

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type Actions struct {
	CanDoFunc func(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool
	DoFunc    func(ctx context.Context, events *model.Events, player *model.Player, req model.ActionRequest, t model.ActionType) (any, error)
}

func (a *Actions) CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool {
	if a.CanDoFunc == nil {
		return false
	}

	return a.CanDoFunc(ctx, events, player, t)
}

func (a *Actions) Do(ctx context.Context, events *model.Events, player *model.Player, req model.ActionRequest, t model.ActionType) (any, error) {
	if a.DoFunc == nil {
		return nil, nil
	}

	return a.DoFunc(ctx, events, player, req, t)
}

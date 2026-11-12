package mocks

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/model"
	"context"
)

type Actions struct {
	CanDoFunc func(ctx context.Context, actor *player.Actor, t model.ActionType) bool
	DoFunc    func(ctx context.Context, actor *player.Actor, req actions.Request, t model.ActionType) (any, error)
}

func (a *Actions) CanDo(ctx context.Context, actor *player.Actor, t model.ActionType) bool {
	if a.CanDoFunc == nil {
		return false
	}

	return a.CanDoFunc(ctx, actor, t)
}

func (a *Actions) Do(ctx context.Context, actor *player.Actor, req actions.Request, t model.ActionType) (any, error) {
	if a.DoFunc == nil {
		return nil, nil
	}

	return a.DoFunc(ctx, actor, req, t)
}

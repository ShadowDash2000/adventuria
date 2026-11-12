package actions

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/model"
	"context"
)

type Request any

type Action interface {
	Type() model.ActionType
	Categories() []string
	InCategory(string) bool
	InCategories(categories []string) bool
	CanDo(ctx context.Context, actor *player.Actor) bool
	Do(ctx context.Context, actor *player.Actor, actionReq Request) (any, error)
}

package board

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type repositoryNotify interface {
	PlayerMovedNotify(ctx context.Context, player *model.Player, moves []*model.MoveResult) error
}

type Notify struct {
	repository repositoryNotify
}

func NewBoardNotify(repository repositoryNotify) *Notify {
	return &Notify{
		repository: repository,
	}
}

func (n *Notify) PlayerMovedNotify(ctx context.Context, player *model.Player, moves []*model.MoveResult) error {
	return n.repository.PlayerMovedNotify(ctx, player, moves)
}

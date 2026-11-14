package board

import (
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/model"
	"context"
)

type repositoryNotify interface {
	ActorMovedNotify(ctx context.Context, actorRef actors.ActorRef, moves []*model.MoveResult) error
}

type Notify struct {
	repository repositoryNotify
}

func NewBoardNotify(repository repositoryNotify) *Notify {
	return &Notify{
		repository: repository,
	}
}

func (n *Notify) ActorMovedNotify(ctx context.Context, actorRef actors.ActorRef, moves []*model.MoveResult) error {
	return n.repository.ActorMovedNotify(ctx, actorRef, moves)
}

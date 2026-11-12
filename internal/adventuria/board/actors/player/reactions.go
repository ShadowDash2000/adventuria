package player

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type CellReached interface {
	OnPlayerReached(ctx context.Context, actor *Actor, reached *model.ReachedContext) error
}

type CellLeft interface {
	OnPlayerLeft(ctx context.Context, actor *Actor) error
}

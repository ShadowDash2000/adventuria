package player

import (
	"adventuria/internal/adventuria/cells"
	"context"
)

type CellReached interface {
	OnPlayerReached(ctx context.Context, actor *Actor, reached *cells.ReachedContext) error
}

type CellLeft interface {
	OnPlayerLeft(ctx context.Context, actor *Actor) error
}

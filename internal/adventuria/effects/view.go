package effects

import (
	"adventuria/internal/adventuria/board/actors/player"
	"context"
)

type View interface {
	GetView(ctx context.Context, actor *player.Actor) (any, error)
}

package behaviors

import (
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/board/actors/npc"
	"context"
)

type Behavior interface {
	Move(ctx context.Context, actor *npc.Actor) (actors.Position, error)
}

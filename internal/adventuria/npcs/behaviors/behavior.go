package behaviors

import (
	"adventuria/internal/adventuria/board/actors/npc"
	"adventuria/internal/adventuria/model"
	"context"
)

type Behavior interface {
	Move(ctx context.Context, actor *npc.Actor) (model.ActorPosition, error)
}

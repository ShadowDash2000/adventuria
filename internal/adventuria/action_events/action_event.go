package action_events

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type ActionEvent interface {
	Init(ctx context.Context, npc *model.NPC, player *model.Player) error
	Dispose(ctx context.Context, npc *model.NPC, player *model.Player) error
}

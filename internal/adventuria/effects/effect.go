package effects

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

type Callback func(ctx context.Context)

type Context struct {
	InvItemID string
	Priority  int
}

type Effect interface {
	Data() *model.EffectInfo
	CanUse(ctx context.Context, actor *player.Actor) bool
	Subscribe(ctx context.Context, actor *player.Actor, effectCtx Context, callback Callback) ([]event.Unsubscribe, error)
}

type EffectPersistent interface {
	Subscribe(ctx context.Context, actor *player.Actor) ([]event.Unsubscribe, error)
}

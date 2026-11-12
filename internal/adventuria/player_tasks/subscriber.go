package player_tasks

import (
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

type PlayerTaskCompleteFunc func(ctx context.Context) error

type PlayerTaskSubscriber interface {
	Subscribe(ctx context.Context, events *model.Events, player *model.Player, complete PlayerTaskCompleteFunc, payload string) ([]event.Unsubscribe, error)
}

package action_events

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type ActionEvent interface {
	Data() *model.ActionEventInfo
	Init(ctx context.Context, player *model.Player) error
	Dispose(ctx context.Context, player *model.Player) error
}

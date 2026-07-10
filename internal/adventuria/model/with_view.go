package model

import "context"

type WithView interface {
	GetView(ctx context.Context, events *Events, player *Player) (any, error)
}

package npcs

import (
	"context"
	"time"
)

type notifyRepo interface {
	NextShiftChangeNotify(ctx context.Context, id string, nextShiftChangeAt time.Time) error
}

type Notify struct {
	notify notifyRepo
}

func NewNPCsNotify(n notifyRepo) *Notify {
	return &Notify{
		notify: n,
	}
}

func (n *Notify) NextShiftChangeNotify(ctx context.Context, id string, nextShiftChangeAt time.Time) error {
	return n.notify.NextShiftChangeNotify(ctx, id, nextShiftChangeAt)
}

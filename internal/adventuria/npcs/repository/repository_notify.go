package repository

import (
	"adventuria/pkg/pbrealtime"
	"adventuria/pkg/pbtransaction"
	"context"
	"time"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
)

type Notify struct {
	pb core.App
}

func NewRepositoryNotify(pb core.App) *Notify {
	return &Notify{pb: pb}
}

func (n *Notify) NextShiftChangeNotify(ctx context.Context, id string, nextShiftChangeAt time.Time) error {
	pbtransaction.AfterCommit(ctx, n.pb, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		_ = pbrealtime.Notify(ctx, n.pb, "next_shift_change", struct {
			NPCId             string    `json:"npc_id"`
			NextShiftChangeAt time.Time `json:"next_shift_change_at"`
		}{
			NPCId:             id,
			NextShiftChangeAt: nextShiftChangeAt,
		})
	})

	return nil
}

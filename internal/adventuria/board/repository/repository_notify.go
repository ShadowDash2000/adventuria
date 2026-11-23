package repository

import (
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/model"
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

type moveNotification struct {
	Type           model.MoveType `json:"type"`
	WorldId        string         `json:"world_id"`
	CellLocalOrder int            `json:"cell_local_order"`
	TotalSteps     int            `json:"total_steps"`
	PrevTotalSteps int            `json:"prev_total_steps"`
}

func (n *Notify) ActorMovedNotify(ctx context.Context, actorRef actors.ActorRef, moves []*model.MoveResult) error {
	var data []moveNotification
	for _, move := range moves {
		data = append(data, moveNotification{
			Type:           move.Type,
			WorldId:        move.World.ID(),
			CellLocalOrder: move.Cell.LocalOrder(),
			TotalSteps:     move.TotalSteps,
			PrevTotalSteps: move.PrevTotalSteps,
		})
	}

	pbtransaction.AfterCommit(ctx, n.pb, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		_ = pbrealtime.Notify(ctx, n.pb, "actor_moved", struct {
			ActorId   string             `json:"actor_id"`
			ActorKind actors.ActorKind   `json:"actor_kind"`
			Moves     []moveNotification `json:"moves"`
		}{
			ActorId:   actorRef.Id,
			ActorKind: actorRef.Kind,
			Moves:     data,
		})
	})

	return nil
}

package repository

import (
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
	Type            model.MoveType `json:"type"`
	WorldId         string         `json:"world_id"`
	CellLocalOrder  int            `json:"cell_local_order"`
	CellGlobalOrder int            `json:"cell_global_order"`
	Steps           int            `json:"steps"`
	TotalSteps      int            `json:"total_steps"`
	PrevTotalSteps  int            `json:"prev_total_steps"`
}

func (n *Notify) PlayerMovedNotify(ctx context.Context, player *model.Player, moves []*model.MoveResult) error {
	var data []moveNotification
	for _, move := range moves {
		data = append(data, moveNotification{
			Type:            move.Type,
			WorldId:         move.CurrentWorld.ID(),
			CellLocalOrder:  move.CurrentCell.LocalOrder(),
			CellGlobalOrder: move.CurrentCell.GlobalOrder(),
			Steps:           move.Steps,
			TotalSteps:      move.TotalSteps,
			PrevTotalSteps:  move.PrevTotalSteps,
		})
	}

	pbtransaction.AfterCommit(ctx, n.pb, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		_ = pbrealtime.Notify(ctx, n.pb, "player_moved", struct {
			PlayerId string             `json:"player_id"`
			Moves    []moveNotification `json:"moves"`
		}{
			PlayerId: player.ID(),
			Moves:    data,
		})
	})

	return nil
}

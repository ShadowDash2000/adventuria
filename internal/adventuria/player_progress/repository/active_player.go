package repository

import (
	"adventuria/internal/adventuria/player_progress"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"

	"github.com/pocketbase/pocketbase/tools/types"
)

type activePlayerRow struct {
	PlayerId        string         `db:"player_id"`
	Name            string         `db:"name"`
	Avatar          string         `db:"avatar"`
	Color           string         `db:"color"`
	IsStreamLive    bool           `db:"is_stream_live"`
	ProgressId      string         `db:"progress_id"`
	ProgressUpdated types.DateTime `db:"progress_updated"`
	CurrentWorld    string         `db:"current_world"`
	CellsPassed     int            `db:"cells_passed"`
}

func (a *activePlayerRow) toActivePlayer() *player_progress.ActivePlayer {
	if a == nil {
		return nil
	}

	return &player_progress.ActivePlayer{
		Id:           a.PlayerId,
		Name:         a.Name,
		Avatar:       pbhelper.BuildFilePath(schema.CollectionPlayers, a.PlayerId, a.Avatar),
		Color:        a.Color,
		IsStreamLive: a.IsStreamLive,
		Progress: player_progress.ActivePlayerProgress{
			Id:           a.ProgressId,
			Updated:      a.ProgressUpdated.Time(),
			CurrentWorld: a.CurrentWorld,
			CellsPassed:  a.CellsPassed,
		},
	}
}

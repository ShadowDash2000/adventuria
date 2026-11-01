package game_state

import (
	"adventuria/internal/adventuria/player_progress"
	"time"
)

type activePlayerView struct {
	Id           string                   `json:"id"`
	Name         string                   `json:"name"`
	Avatar       string                   `json:"avatar"`
	Color        string                   `json:"color"`
	IsStreamLive bool                     `json:"is_stream_live"`
	Progress     activePlayerProgressView `json:"progress"`
}

type activePlayerProgressView struct {
	Id           string    `json:"id"`
	Updated      time.Time `json:"updated"`
	CurrentWorld string    `json:"current_world"`
	CellsPassed  int       `json:"cells_passed"`
}

func activePlayerToView(player *player_progress.ActivePlayer) *activePlayerView {
	return &activePlayerView{
		Id:           player.Id,
		Name:         player.Name,
		Avatar:       player.Avatar,
		Color:        player.Color,
		IsStreamLive: player.IsStreamLive,
		Progress: activePlayerProgressView{
			Id:           player.Progress.Id,
			Updated:      player.Progress.Updated,
			CurrentWorld: player.Progress.CurrentWorld,
			CellsPassed:  player.Progress.CellsPassed,
		},
	}
}

func activePlayersToView(players []*player_progress.ActivePlayer) []*activePlayerView {
	res := make([]*activePlayerView, len(players))
	for i, player := range players {
		res[i] = activePlayerToView(player)
	}
	return res
}

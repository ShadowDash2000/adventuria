package player_progress

import "time"

type ActivePlayer struct {
	Id           string
	Name         string
	Avatar       string
	Color        string
	IsStreamLive bool
	Progress     ActivePlayerProgress
}

type ActivePlayerProgress struct {
	Id           string
	Updated      time.Time
	CurrentWorld string
	CellsPassed  int
}

package games

import "github.com/pocketbase/pocketbase/core"

type GameType struct {
	IdDb     string
	Name     string
	Checksum string
}

type GameTypeRecord interface {
	core.RecordProxy
	UpdatableRecord

	ID() string
	SetIdDb(string)
	Name() string
	SetName(string)
}

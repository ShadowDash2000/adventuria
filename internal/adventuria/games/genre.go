package games

import "github.com/pocketbase/pocketbase/core"

type Genre struct {
	IdDb     string
	Name     string
	Checksum string
}

type GenreRecord interface {
	core.RecordProxy
	UpdatableRecord

	ID() string
	SetIdDb(string)
	Name() string
	SetName(string)
}

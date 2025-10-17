package games

import "github.com/pocketbase/pocketbase/core"

type Genre struct {
	IdDb     uint64
	Name     string
	Checksum string
}

type GenreRecord interface {
	core.RecordProxy
	UpdatableRecord

	ID() string
	SetIdDb(uint64)
	Name() string
	SetName(string)
}

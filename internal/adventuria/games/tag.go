package games

import "github.com/pocketbase/pocketbase/core"

type Tag struct {
	IdDb     uint64
	Name     string
	Checksum string
}

type TagRecord interface {
	core.RecordProxy
	UpdatableRecord

	ID() string
	SetIdDb(uint64)
	Name() string
	SetName(string)
}

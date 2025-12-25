package games

import "github.com/pocketbase/pocketbase/core"

type Tag struct {
	IdDb     string
	Name     string
	Checksum string
}

type TagRecord interface {
	core.RecordProxy
	UpdatableRecord

	ID() string
	SetIdDb(string)
	Name() string
	SetName(string)
}

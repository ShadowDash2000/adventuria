package games

import (
	"github.com/pocketbase/pocketbase/core"
)

type Platform struct {
	IdDb     string
	Name     string
	Checksum string
}

type PlatformRecord interface {
	core.RecordProxy
	UpdatableRecord

	ID() string
	IdDb() string
	SetIdDb(string)
	Name() string
	SetName(string)
}

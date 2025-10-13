package games

import (
	"github.com/pocketbase/pocketbase/core"
)

type Platform struct {
	IdDb     uint64
	Name     string
	Checksum string
}

type PlatformRecord interface {
	core.RecordProxy
	UpdatableRecord

	ID() string
	IdDb() uint64
	SetIdDb(uint64)
	Name() string
	SetName(string)
}

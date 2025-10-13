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

	ID() string
	IdDb() uint64
	SetIdDb(int)
	Name() string
	SetName(string)

	Checksum() string
	SetChecksum(string)
}

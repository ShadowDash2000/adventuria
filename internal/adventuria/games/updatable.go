package games

import "github.com/pocketbase/pocketbase/core"

type UpdatableRecord interface {
	core.RecordProxy

	IdDb() uint64
	Checksum() string
	SetChecksum(string)
}

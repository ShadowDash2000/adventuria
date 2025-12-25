package games

import "github.com/pocketbase/pocketbase/core"

type UpdatableRecord interface {
	core.RecordProxy

	IdDb() string
	Checksum() string
	SetChecksum(string)
}

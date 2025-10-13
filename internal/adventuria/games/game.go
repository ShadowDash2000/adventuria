package games

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Game struct {
	IdDb        uint64
	Name        string
	ReleaseDate types.DateTime
	Platforms   []uint64
	Checksum    string
}

type GameRecord interface {
	core.RecordProxy

	ID() string
	IdDb() uint64
	SetIdDb(int)
	Name() string
	SetName(string)
	ReleaseDate() types.DateTime
	SetReleaseDate(types.DateTime)
	Platforms() []string
	SetPlatforms([]string)

	Checksum() string
	SetChecksum(string)
}

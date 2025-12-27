package games

import (
	"adventuria/internal/adventuria"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Game struct {
	IdDb        string
	Name        string
	Slug        string
	ReleaseDate types.DateTime
	Platforms   CollectionReference
	Developers  CollectionReference
	Publishers  CollectionReference
	Genres      CollectionReference
	Tags        CollectionReference
	Themes      CollectionReference
	SteamAppId  uint64
	Cover       string
	Checksum    string
}

type GameRecord interface {
	adventuria.ActivityRecord
	UpdatableRecord
}

func NewGameFromRecord(record *core.Record) GameRecord {
	a := &adventuria.ActivityRecordBase{}
	a.SetProxyRecord(record)
	return a
}

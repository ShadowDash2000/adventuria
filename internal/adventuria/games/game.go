package games

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Game struct {
	IdDb        uint64
	Name        string
	Slug        string
	ReleaseDate types.DateTime
	Platforms   CollectionReference
	Developers  CollectionReference
	Publishers  CollectionReference
	Genres      CollectionReference
	Tags        CollectionReference
	SteamAppId  uint64
	Cover       string
	Checksum    string
}

type GameRecord interface {
	core.RecordProxy
	UpdatableRecord

	ID() string
	IdDb() uint64
	SetIdDb(uint64)
	Name() string
	SetName(string)
	Slug() string
	SetSlug(string)
	ReleaseDate() types.DateTime
	SetReleaseDate(types.DateTime)
	Platforms() []string
	SetPlatforms([]string)
	Developers() []string
	SetDevelopers([]string)
	Publishers() []string
	SetPublishers([]string)
	Genres() []string
	SetGenres([]string)
	Tags() []string
	SetTags([]string)
	SteamAppId() uint64
	SetSteamAppId(uint64)
	SteamAppPrice() uint
	SetSteamAppPrice(uint)
	HltbId() int
	SetHltbId(int)
	Campaign() float64
	SetCampaign(float64)
	Cover() string
	SetCover(string)
}

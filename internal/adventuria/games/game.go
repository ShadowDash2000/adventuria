package games

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

// TODO maybe we should combine games with movies in one DB collection

type Game struct {
	IdDb        uint64
	Name        string
	ReleaseDate types.DateTime
	Platforms   CollectionReference
	Companies   CollectionReference
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
	ReleaseDate() types.DateTime
	SetReleaseDate(types.DateTime)
	Platforms() []string
	SetPlatforms([]string)
	Companies() []string
	SetCompanies([]string)
	SteamAppId() uint64
	SetSteamAppId(uint64)
	SteamAppPrice() int
	SetSteamAppPrice(int)
	CampaignTime() int
	SetCampaignTime(int)
	Cover() string
	SetCover(string)
}

package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type ActivityRecord interface {
	core.RecordProxy

	ID() string
	IdDb() string
	SetIdDb(string)
	Type() ActivityType
	SetType(ActivityType)
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
	Themes() []string
	SetThemes([]string)
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

type ActivityType string

const (
	ActivityTypeGame  ActivityType = "game"
	ActivityTypeMovie ActivityType = "movie"
	ActivityTypeGym   ActivityType = "gym"
)

var ActivityTypes = map[ActivityType]bool{
	ActivityTypeGame:  true,
	ActivityTypeMovie: true,
	ActivityTypeGym:   true,
}

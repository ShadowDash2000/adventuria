package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type GameFilterRecord interface {
	core.RecordProxy

	ID() string
	Name() string
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
	MinPrice() int
	SetMinPrice(int)
	MaxPrice() int
	SetMaxPrice(int)
	ReleaseDateFrom() types.DateTime
	SetReleaseDateFrom(types.DateTime)
	ReleaseDateTo() types.DateTime
	SetReleaseDateTo(types.DateTime)
	MinCampaignTime() float64
	SetMinCampaignTime(float64)
	MaxCampaignTime() float64
	SetMaxCampaignTime(float64)
	Games() []string
	SetGames([]string)
}

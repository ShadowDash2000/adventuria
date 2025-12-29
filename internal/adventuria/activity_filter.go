package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type ActivityFilterRecord interface {
	core.RecordProxy

	ID() string
	Type() ActivityType
	SetType(ActivityType)
	Name() string
	Platforms() []string
	SetPlatforms([]string)
	PlatformsStrict() bool
	SetPlatformsStrict(bool)
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
	Activities() []string
	SetActivities([]string)
}

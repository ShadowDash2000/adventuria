package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type ActivityFilterBase struct {
	core.BaseRecordProxy
}

func NewActivityFilterFromRecord(record *core.Record) ActivityFilterRecord {
	f := &ActivityFilterBase{}
	f.SetProxyRecord(record)
	return f
}

func (a *ActivityFilterBase) ID() string {
	return a.Id
}

func (a *ActivityFilterBase) Type() ActivityType {
	return ActivityType(a.GetString("type"))
}

func (a *ActivityFilterBase) SetType(t ActivityType) {
	a.Set("type", t)
}

func (a *ActivityFilterBase) Name() string {
	return a.GetString("name")
}

func (a *ActivityFilterBase) Platforms() []string {
	return a.GetStringSlice("platforms")
}

func (a *ActivityFilterBase) SetPlatforms(platforms []string) {
	a.Set("platforms", platforms)
}

func (a *ActivityFilterBase) PlatformsStrict() bool {
	return a.GetBool("platforms_strict")
}

func (a *ActivityFilterBase) SetPlatformsStrict(strict bool) {
	a.Set("platforms_strict", strict)
}

func (a *ActivityFilterBase) Developers() []string {
	return a.GetStringSlice("developers")
}

func (a *ActivityFilterBase) SetDevelopers(developers []string) {
	a.Set("developers", developers)
}

func (a *ActivityFilterBase) Publishers() []string {
	return a.GetStringSlice("publishers")
}

func (a *ActivityFilterBase) SetPublishers(publishers []string) {
	a.Set("publishers", publishers)
}

func (a *ActivityFilterBase) Genres() []string {
	return a.GetStringSlice("genres")
}

func (a *ActivityFilterBase) SetGenres(genres []string) {
	a.Set("genres", genres)
}

func (a *ActivityFilterBase) Tags() []string {
	return a.GetStringSlice("tags")
}

func (a *ActivityFilterBase) SetTags(tags []string) {
	a.Set("tags", tags)
}

func (a *ActivityFilterBase) MinPrice() int {
	return a.GetInt("min_price")
}

func (a *ActivityFilterBase) SetMinPrice(minPrice int) {
	a.Set("min_price", minPrice)
}

func (a *ActivityFilterBase) MaxPrice() int {
	return a.GetInt("max_price")
}

func (a *ActivityFilterBase) SetMaxPrice(maxPrice int) {
	a.Set("max_price", maxPrice)
}

func (a *ActivityFilterBase) ReleaseDateFrom() types.DateTime {
	return a.GetDateTime("release_date_from")
}

func (a *ActivityFilterBase) SetReleaseDateFrom(releaseDateFrom types.DateTime) {
	a.Set("release_date_from", releaseDateFrom)
}

func (a *ActivityFilterBase) ReleaseDateTo() types.DateTime {
	return a.GetDateTime("release_date_to")
}

func (a *ActivityFilterBase) SetReleaseDateTo(releaseDateTo types.DateTime) {
	a.Set("release_date_to", releaseDateTo)
}

func (a *ActivityFilterBase) MinCampaignTime() float64 {
	return a.GetFloat("min_campaign_time")
}

func (a *ActivityFilterBase) SetMinCampaignTime(minCampaignTime float64) {
	a.Set("min_campaign_time", minCampaignTime)
}

func (a *ActivityFilterBase) MaxCampaignTime() float64 {
	return a.GetFloat("max_campaign_time")
}

func (a *ActivityFilterBase) SetMaxCampaignTime(maxCampaignTime float64) {
	a.Set("max_campaign_time", maxCampaignTime)
}

func (a *ActivityFilterBase) Games() []string {
	return a.GetStringSlice("games")
}

func (a *ActivityFilterBase) SetGames(games []string) {
	a.Set("games", games)
}

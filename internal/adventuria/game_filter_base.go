package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type GameFilterBase struct {
	core.BaseRecordProxy
}

func NewGameFilterFromRecord(record *core.Record) GameFilterRecord {
	return &GameFilterBase{}
}

func (g *GameFilterBase) ID() string {
	return g.Id
}

func (g *GameFilterBase) Name() string {
	return g.GetString("name")
}

func (g *GameFilterBase) Platforms() []string {
	return g.GetStringSlice("platforms")
}

func (g *GameFilterBase) SetPlatforms(platforms []string) {
	g.Set("platforms", platforms)
}

func (g *GameFilterBase) Developers() []string {
	return g.GetStringSlice("developers")
}

func (g *GameFilterBase) SetDevelopers(developers []string) {
	g.Set("developers", developers)
}

func (g *GameFilterBase) Publishers() []string {
	return g.GetStringSlice("publishers")
}

func (g *GameFilterBase) SetPublishers(publishers []string) {
	g.Set("publishers", publishers)
}

func (g *GameFilterBase) Genres() []string {
	return g.GetStringSlice("genres")
}

func (g *GameFilterBase) SetGenres(genres []string) {
	g.Set("genres", genres)
}

func (g *GameFilterBase) Tags() []string {
	return g.GetStringSlice("tags")
}

func (g *GameFilterBase) SetTags(tags []string) {
	g.Set("tags", tags)
}

func (g *GameFilterBase) MinPrice() int {
	return g.GetInt("min_price")
}

func (g *GameFilterBase) SetMinPrice(minPrice int) {
	g.Set("min_price", minPrice)
}

func (g *GameFilterBase) MaxPrice() int {
	return g.GetInt("max_price")
}

func (g *GameFilterBase) SetMaxPrice(maxPrice int) {
	g.Set("max_price", maxPrice)
}

func (g *GameFilterBase) ReleaseDateFrom() types.DateTime {
	return g.GetDateTime("release_date_from")
}

func (g *GameFilterBase) SetReleaseDateFrom(releaseDateFrom types.DateTime) {
	g.Set("release_date_from", releaseDateFrom)
}

func (g *GameFilterBase) ReleaseDateTo() types.DateTime {
	return g.GetDateTime("release_date_to")
}

func (g *GameFilterBase) SetReleaseDateTo(releaseDateTo types.DateTime) {
	g.Set("release_date_to", releaseDateTo)
}

func (g *GameFilterBase) MinCampaignTime() float64 {
	return g.GetFloat("min_campaign_time")
}

func (g *GameFilterBase) SetMinCampaignTime(minCampaignTime float64) {
	g.Set("min_campaign_time", minCampaignTime)
}

func (g *GameFilterBase) MaxCampaignTime() float64 {
	return g.GetFloat("max_campaign_time")
}

func (g *GameFilterBase) SetMaxCampaignTime(maxCampaignTime float64) {
	g.Set("max_campaign_time", maxCampaignTime)
}

func (g *GameFilterBase) Games() []string {
	return g.GetStringSlice("games")
}

func (g *GameFilterBase) SetGames(games []string) {
	g.Set("games", games)
}

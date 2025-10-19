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

func (g GameFilterBase) ID() string {
	return g.Id
}

func (g GameFilterBase) Name() string {
	return g.GetString("name")
}

func (g GameFilterBase) Platforms() []string {
	return g.GetStringSlice("platforms")
}

func (g GameFilterBase) Developers() []string {
	return g.GetStringSlice("developers")
}

func (g GameFilterBase) Publishers() []string {
	return g.GetStringSlice("publishers")
}

func (g GameFilterBase) Genres() []string {
	return g.GetStringSlice("genres")
}

func (g GameFilterBase) Tags() []string {
	return g.GetStringSlice("tags")
}

func (g GameFilterBase) MinPrice() int {
	return g.GetInt("min_price")
}

func (g GameFilterBase) MaxPrice() int {
	return g.GetInt("max_price")
}

func (g GameFilterBase) ReleaseDateFrom() types.DateTime {
	return g.GetDateTime("release_date_from")
}

func (g GameFilterBase) ReleaseDateTo() types.DateTime {
	return g.GetDateTime("release_date_to")
}

func (g GameFilterBase) MinCampaignTime() float64 {
	return g.GetFloat("min_campaign_time")
}

func (g GameFilterBase) MaxCampaignTime() float64 {
	return g.GetFloat("max_campaign_time")
}

func (g GameFilterBase) Games() []string {
	return g.GetStringSlice("games")
}

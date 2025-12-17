package games

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type GameRecordBase struct {
	core.BaseRecordProxy
}

func NewGameFromRecord(record *core.Record) GameRecord {
	g := &GameRecordBase{}
	g.SetProxyRecord(record)
	return g
}

func (g *GameRecordBase) ID() string {
	return g.Id
}

func (g *GameRecordBase) IdDb() uint64 {
	return uint64(g.GetInt("id_db"))
}

func (g *GameRecordBase) SetIdDb(id uint64) {
	g.Set("id_db", int(id))
}

func (g *GameRecordBase) Name() string {
	return g.GetString("name")
}

func (g *GameRecordBase) SetName(name string) {
	g.Set("name", name)
}

func (g *GameRecordBase) Slug() string {
	return g.GetString("slug")
}

func (g *GameRecordBase) SetSlug(slug string) {
	g.Set("slug", slug)
}

func (g *GameRecordBase) ReleaseDate() types.DateTime {
	return g.GetDateTime("releaseDate")
}

func (g *GameRecordBase) SetReleaseDate(date types.DateTime) {
	g.Set("release_date", date)
}

func (g *GameRecordBase) Platforms() []string {
	return g.GetStringSlice("platforms")
}

func (g *GameRecordBase) SetPlatforms(ids []string) {
	g.Set("platforms", ids)
}

func (g *GameRecordBase) Developers() []string {
	return g.GetStringSlice("developers")
}

func (g *GameRecordBase) SetDevelopers(ids []string) {
	g.Set("developers", ids)
}

func (g *GameRecordBase) Publishers() []string {
	return g.GetStringSlice("publishers")
}

func (g *GameRecordBase) SetPublishers(ids []string) {
	g.Set("publishers", ids)
}

func (g *GameRecordBase) Genres() []string {
	return g.GetStringSlice("genres")
}

func (g *GameRecordBase) SetGenres(ids []string) {
	g.Set("genres", ids)
}

func (g *GameRecordBase) Tags() []string {
	return g.GetStringSlice("tags")
}

func (g *GameRecordBase) SetTags(ids []string) {
	g.Set("tags", ids)
}

func (g *GameRecordBase) SteamAppId() uint64 {
	return uint64(g.GetInt("steam_app_id"))
}

func (g *GameRecordBase) SetSteamAppId(id uint64) {
	g.Set("steam_app_id", int(id))
}

func (g *GameRecordBase) SteamAppPrice() uint {
	return uint(g.GetInt("steam_app_price"))
}

func (g *GameRecordBase) SetSteamAppPrice(price uint) {
	g.Set("steam_app_price", price)
}

func (g *GameRecordBase) SteamSpy() string {
	return g.GetString("steam_spy")
}

func (g *GameRecordBase) SetSteamSpy(id string) {
	g.Set("steam_spy", id)
}

func (g *GameRecordBase) HltbId() int {
	return g.GetInt("hltb_id")
}

func (g *GameRecordBase) SetHltbId(id int) {
	g.Set("hltb_id", id)
}

func (g *GameRecordBase) Campaign() float64 {
	return g.GetFloat("hltb_campaign_time")
}

func (g *GameRecordBase) SetCampaign(campaign float64) {
	g.Set("hltb_campaign_time", campaign)
}

func (g *GameRecordBase) Cover() string {
	return g.GetString("cover")
}

func (g *GameRecordBase) SetCover(url string) {
	g.Set("cover", url)
}

func (g *GameRecordBase) Checksum() string {
	return g.GetString("checksum")
}

func (g *GameRecordBase) SetChecksum(checksum string) {
	g.Set("checksum", checksum)
}

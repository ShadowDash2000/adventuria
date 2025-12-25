package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type ActivityRecordBase struct {
	core.BaseRecordProxy
}

func NewActivityFromRecord(record *core.Record) ActivityRecord {
	a := &ActivityRecordBase{}
	a.SetProxyRecord(record)
	return a
}

func (a *ActivityRecordBase) ID() string {
	return a.Id
}

func (a *ActivityRecordBase) IdDb() string {
	return a.GetString("id_db")
}

func (a *ActivityRecordBase) SetIdDb(id string) {
	a.Set("id_db", id)
}

func (a *ActivityRecordBase) Type() ActivityType {
	return ActivityType(a.GetString("type"))
}

func (a *ActivityRecordBase) SetType(t ActivityType) {
	a.Set("type", t)
}

func (a *ActivityRecordBase) Name() string {
	return a.GetString("name")
}

func (a *ActivityRecordBase) SetName(name string) {
	a.Set("name", name)
}

func (a *ActivityRecordBase) Slug() string {
	return a.GetString("slug")
}

func (a *ActivityRecordBase) SetSlug(slug string) {
	a.Set("slug", slug)
}

func (a *ActivityRecordBase) ReleaseDate() types.DateTime {
	return a.GetDateTime("releaseDate")
}

func (a *ActivityRecordBase) SetReleaseDate(date types.DateTime) {
	a.Set("release_date", date)
}

func (a *ActivityRecordBase) Platforms() []string {
	return a.GetStringSlice("platforms")
}

func (a *ActivityRecordBase) SetPlatforms(ids []string) {
	a.Set("platforms", ids)
}

func (a *ActivityRecordBase) Developers() []string {
	return a.GetStringSlice("developers")
}

func (a *ActivityRecordBase) SetDevelopers(ids []string) {
	a.Set("developers", ids)
}

func (a *ActivityRecordBase) Publishers() []string {
	return a.GetStringSlice("publishers")
}

func (a *ActivityRecordBase) SetPublishers(ids []string) {
	a.Set("publishers", ids)
}

func (a *ActivityRecordBase) Genres() []string {
	return a.GetStringSlice("genres")
}

func (a *ActivityRecordBase) SetGenres(ids []string) {
	a.Set("genres", ids)
}

func (a *ActivityRecordBase) Tags() []string {
	return a.GetStringSlice("tags")
}

func (a *ActivityRecordBase) SetTags(ids []string) {
	a.Set("tags", ids)
}

func (a *ActivityRecordBase) SteamAppId() uint64 {
	return uint64(a.GetInt("steam_app_id"))
}

func (a *ActivityRecordBase) SetSteamAppId(id uint64) {
	a.Set("steam_app_id", int(id))
}

func (a *ActivityRecordBase) SteamAppPrice() uint {
	return uint(a.GetInt("steam_app_price"))
}

func (a *ActivityRecordBase) SetSteamAppPrice(price uint) {
	a.Set("steam_app_price", price)
}

func (a *ActivityRecordBase) HltbId() int {
	return a.GetInt("hltb_id")
}

func (a *ActivityRecordBase) SetHltbId(id int) {
	a.Set("hltb_id", id)
}

func (a *ActivityRecordBase) Campaign() float64 {
	return a.GetFloat("hltb_campaign_time")
}

func (a *ActivityRecordBase) SetCampaign(campaign float64) {
	a.Set("hltb_campaign_time", campaign)
}

func (a *ActivityRecordBase) Cover() string {
	return a.GetString("cover")
}

func (a *ActivityRecordBase) SetCover(url string) {
	a.Set("cover", url)
}

func (a *ActivityRecordBase) Checksum() string {
	return a.GetString("checksum")
}

func (a *ActivityRecordBase) SetChecksum(checksum string) {
	a.Set("checksum", checksum)
}

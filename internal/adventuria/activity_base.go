package adventuria

import (
	"adventuria/internal/adventuria/schema"

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
	return a.GetString(schema.ActivitySchema.IdDb)
}

func (a *ActivityRecordBase) SetIdDb(id string) {
	a.Set(schema.ActivitySchema.IdDb, id)
}

func (a *ActivityRecordBase) Type() ActivityType {
	return ActivityType(a.GetString(schema.ActivitySchema.Type))
}

func (a *ActivityRecordBase) SetType(t ActivityType) {
	a.Set(schema.ActivitySchema.Type, t)
}

func (a *ActivityRecordBase) Name() string {
	return a.GetString(schema.ActivitySchema.Name)
}

func (a *ActivityRecordBase) SetName(name string) {
	a.Set(schema.ActivitySchema.Name, name)
}

func (a *ActivityRecordBase) Slug() string {
	return a.GetString(schema.ActivitySchema.Slug)
}

func (a *ActivityRecordBase) SetSlug(slug string) {
	a.Set(schema.ActivitySchema.Slug, slug)
}

func (a *ActivityRecordBase) GameType() string {
	return a.GetString(schema.ActivitySchema.GameType)
}

func (a *ActivityRecordBase) SetGameType(gameType string) {
	a.Set(schema.ActivitySchema.GameType, gameType)
}

func (a *ActivityRecordBase) ReleaseDate() types.DateTime {
	return a.GetDateTime(schema.ActivitySchema.ReleaseDate)
}

func (a *ActivityRecordBase) SetReleaseDate(date types.DateTime) {
	a.Set(schema.ActivitySchema.ReleaseDate, date)
}

func (a *ActivityRecordBase) Platforms() []string {
	return a.GetStringSlice(schema.ActivitySchema.Platforms)
}

func (a *ActivityRecordBase) SetPlatforms(ids []string) {
	a.Set(schema.ActivitySchema.Platforms, ids)
}

func (a *ActivityRecordBase) Developers() []string {
	return a.GetStringSlice(schema.ActivitySchema.Developers)
}

func (a *ActivityRecordBase) SetDevelopers(ids []string) {
	a.Set(schema.ActivitySchema.Developers, ids)
}

func (a *ActivityRecordBase) Publishers() []string {
	return a.GetStringSlice(schema.ActivitySchema.Publishers)
}

func (a *ActivityRecordBase) SetPublishers(ids []string) {
	a.Set(schema.ActivitySchema.Publishers, ids)
}

func (a *ActivityRecordBase) Genres() []string {
	return a.GetStringSlice(schema.ActivitySchema.Genres)
}

func (a *ActivityRecordBase) SetGenres(ids []string) {
	a.Set(schema.ActivitySchema.Genres, ids)
}

func (a *ActivityRecordBase) Tags() []string {
	return a.GetStringSlice(schema.ActivitySchema.Tags)
}

func (a *ActivityRecordBase) SetTags(ids []string) {
	a.Set(schema.ActivitySchema.Tags, ids)
}

func (a *ActivityRecordBase) Themes() []string {
	return a.GetStringSlice(schema.ActivitySchema.Themes)
}

func (a *ActivityRecordBase) SetThemes(ids []string) {
	a.Set(schema.ActivitySchema.Themes, ids)
}

func (a *ActivityRecordBase) SteamAppId() uint64 {
	return uint64(a.GetInt(schema.ActivitySchema.SteamAppId))
}

func (a *ActivityRecordBase) SetSteamAppId(id uint64) {
	a.Set(schema.ActivitySchema.SteamAppId, int(id))
}

func (a *ActivityRecordBase) SteamAppPrice() uint {
	return uint(a.GetInt(schema.ActivitySchema.SteamAppPrice))
}

func (a *ActivityRecordBase) SetSteamAppPrice(price uint) {
	a.Set(schema.ActivitySchema.SteamAppPrice, price)
}

func (a *ActivityRecordBase) HltbId() int {
	return a.GetInt(schema.ActivitySchema.HltbId)
}

func (a *ActivityRecordBase) SetHltbId(id int) {
	a.Set(schema.ActivitySchema.HltbId, id)
}

func (a *ActivityRecordBase) Campaign() float64 {
	return a.GetFloat(schema.ActivitySchema.HltbCampaignTime)
}

func (a *ActivityRecordBase) SetCampaign(campaign float64) {
	a.Set(schema.ActivitySchema.HltbCampaignTime, campaign)
}

func (a *ActivityRecordBase) Cover() string {
	return a.GetString(schema.ActivitySchema.Cover)
}

func (a *ActivityRecordBase) SetCover(url string) {
	a.Set(schema.ActivitySchema.Cover, url)
}

func (a *ActivityRecordBase) Checksum() string {
	return a.GetString(schema.ActivitySchema.Checksum)
}

func (a *ActivityRecordBase) SetChecksum(checksum string) {
	a.Set(schema.ActivitySchema.Checksum, checksum)
}

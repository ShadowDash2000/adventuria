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

func (c *GameRecordBase) ID() string {
	return c.Id
}

func (c *GameRecordBase) IdDb() uint64 {
	return uint64(c.GetInt("id_db"))
}

func (c *GameRecordBase) SetIdDb(id uint64) {
	c.Set("id_db", int(id))
}

func (c *GameRecordBase) Name() string {
	return c.GetString("name")
}

func (c *GameRecordBase) SetName(name string) {
	c.Set("name", name)
}

func (c *GameRecordBase) ReleaseDate() types.DateTime {
	return c.GetDateTime("releaseDate")
}

func (c *GameRecordBase) SetReleaseDate(date types.DateTime) {
	c.Set("release_date", date)
}

func (c *GameRecordBase) Platforms() []string {
	return c.GetStringSlice("platforms")
}

func (c *GameRecordBase) SetPlatforms(ids []string) {
	c.Set("platforms", ids)
}

func (c *GameRecordBase) Companies() []string {
	return c.GetStringSlice("companies")
}

func (c *GameRecordBase) SetCompanies(ids []string) {
	c.Set("companies", ids)
}

func (c *GameRecordBase) SteamAppId() uint64 {
	return uint64(c.GetInt("steam_app_id"))
}

func (c *GameRecordBase) SetSteamAppId(id uint64) {
	c.Set("steam_app_id", int(id))
}

func (c *GameRecordBase) Cover() string {
	return c.GetString("cover")
}

func (c *GameRecordBase) SetCover(url string) {
	c.Set("cover", url)
}

func (c *GameRecordBase) Checksum() string {
	return c.GetString("checksum")
}

func (c *GameRecordBase) SetChecksum(checksum string) {
	c.Set("checksum", checksum)
}

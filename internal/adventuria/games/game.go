package games

import (
	"github.com/bestnite/go-igdb/proto"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Game struct {
	core.BaseRecordProxy
}

func (g *Game) IGDBID() int {
	return g.GetInt("igdb_id")
}

func (g *Game) SetIGDBID(id int) {
	g.Set("igdb_id", id)
}

func (g *Game) Name() string {
	return g.GetString("name")
}

func (g *Game) SetName(name string) {
	g.Set("name", name)
}

func (g *Game) FirstReleaseDate() types.DateTime {
	return g.GetDateTime("firstReleaseDate")
}

func (g *Game) SetFirstReleaseDate(date types.DateTime) {
	g.Set("firstReleaseDate", date)
}

func (g *Game) Platforms() []string {
	return g.GetStringSlice("platforms")
}

func (g *Game) SetPlatforms(platforms []string) {
	g.Set("platforms", platforms)
}

func (g *Game) Checksum() string {
	return g.GetString("checksum")
}

func (g *Game) SetChecksum(checksum string) {
	g.Set("checksum", checksum)
}

func (g *Game) Data() *proto.Platform {
	var data *proto.Platform
	g.UnmarshalJSONField("data", &data)
	return data
}

func (g *Game) SetData(data *proto.Game) {
	g.Set("data", data)
}

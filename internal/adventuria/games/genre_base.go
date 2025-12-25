package games

import "github.com/pocketbase/pocketbase/core"

type GenreRecordBase struct {
	core.BaseRecordProxy
}

func NewGenreFromRecord(record *core.Record) GenreRecord {
	g := &GenreRecordBase{}
	g.SetProxyRecord(record)
	return g
}

func (g *GenreRecordBase) ID() string {
	return g.Id
}

func (g *GenreRecordBase) IdDb() string {
	return g.GetString("id_db")
}

func (g *GenreRecordBase) SetIdDb(id string) {
	g.Set("id_db", id)
}

func (g *GenreRecordBase) Name() string {
	return g.GetString("name")
}

func (g *GenreRecordBase) SetName(name string) {
	g.Set("name", name)
}

func (g *GenreRecordBase) Checksum() string {
	return g.GetString("checksum")
}

func (g *GenreRecordBase) SetChecksum(checksum string) {
	g.Set("checksum", checksum)
}

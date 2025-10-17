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

func (g *GenreRecordBase) IdDb() uint64 {
	return uint64(g.GetInt("id_db"))
}

func (g *GenreRecordBase) SetIdDb(id uint64) {
	g.Set("id_db", int(id))
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

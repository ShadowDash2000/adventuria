package games

import "github.com/pocketbase/pocketbase/core"

type GameTypeRecordBase struct {
	core.BaseRecordProxy
}

func NewGameTypeFromRecord(record *core.Record) GameTypeRecord {
	t := &GameTypeRecordBase{}
	t.SetProxyRecord(record)
	return t
}

func (t *GameTypeRecordBase) ID() string {
	return t.Id
}

func (t *GameTypeRecordBase) IdDb() string {
	return t.GetString("id_db")
}

func (t *GameTypeRecordBase) SetIdDb(id string) {
	t.Set("id_db", id)
}

func (t *GameTypeRecordBase) Name() string {
	return t.GetString("name")
}

func (t *GameTypeRecordBase) SetName(name string) {
	t.Set("name", name)
}

func (t *GameTypeRecordBase) Checksum() string {
	return t.GetString("checksum")
}

func (t *GameTypeRecordBase) SetChecksum(checksum string) {
	t.Set("checksum", checksum)
}

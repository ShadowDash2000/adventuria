package games

import "github.com/pocketbase/pocketbase/core"

type TagRecordBase struct {
	core.BaseRecordProxy
}

func NewTagFromRecord(record *core.Record) TagRecord {
	t := &TagRecordBase{}
	t.SetProxyRecord(record)
	return t
}

func (t *TagRecordBase) ID() string {
	return t.Id
}

func (t *TagRecordBase) IdDb() string {
	return t.GetString("id_db")
}

func (t *TagRecordBase) SetIdDb(id string) {
	t.Set("id_db", id)
}

func (t *TagRecordBase) Name() string {
	return t.GetString("name")
}

func (t *TagRecordBase) SetName(name string) {
	t.Set("name", name)
}

func (t *TagRecordBase) Checksum() string {
	return t.GetString("checksum")
}

func (t *TagRecordBase) SetChecksum(checksum string) {
	t.Set("checksum", checksum)
}

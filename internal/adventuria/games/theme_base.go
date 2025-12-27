package games

import "github.com/pocketbase/pocketbase/core"

type ThemeRecordBase struct {
	core.BaseRecordProxy
}

func NewThemeFromRecord(record *core.Record) ThemeRecord {
	t := &ThemeRecordBase{}
	t.SetProxyRecord(record)
	return t
}

func (t *ThemeRecordBase) ID() string {
	return t.Id
}

func (t *ThemeRecordBase) IdDb() string {
	return t.GetString("id_db")
}

func (t *ThemeRecordBase) SetIdDb(id string) {
	t.Set("id_db", id)
}

func (t *ThemeRecordBase) Name() string {
	return t.GetString("name")
}

func (t *ThemeRecordBase) SetName(name string) {
	t.Set("name", name)
}

func (t *ThemeRecordBase) Checksum() string {
	return t.GetString("checksum")
}

func (t *ThemeRecordBase) SetChecksum(checksum string) {
	t.Set("checksum", checksum)
}

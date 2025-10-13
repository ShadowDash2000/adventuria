package games

import "github.com/pocketbase/pocketbase/core"

type PlatformRecordBase struct {
	core.BaseRecordProxy
}

func NewPlatformFromRecord(record *core.Record) PlatformRecord {
	p := &PlatformRecordBase{}
	p.SetProxyRecord(record)
	return p
}

func (c *PlatformRecordBase) ID() string {
	return c.Id
}

func (c *PlatformRecordBase) IdDb() uint64 {
	return uint64(c.GetInt("id_db"))
}

func (c *PlatformRecordBase) SetIdDb(id uint64) {
	c.Set("id_db", int(id))
}

func (c *PlatformRecordBase) Name() string {
	return c.GetString("name")
}

func (c *PlatformRecordBase) SetName(name string) {
	c.Set("name", name)
}

func (c *PlatformRecordBase) Checksum() string {
	return c.GetString("checksum")
}

func (c *PlatformRecordBase) SetChecksum(checksum string) {
	c.Set("checksum", checksum)
}

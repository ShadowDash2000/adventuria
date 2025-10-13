package games

import "github.com/pocketbase/pocketbase/core"

type CompanyRecordBase struct {
	core.BaseRecordProxy
}

func NewCompanyFromRecord(record *core.Record) CompanyRecord {
	c := &CompanyRecordBase{}
	c.SetProxyRecord(record)
	return c
}

func (c *CompanyRecordBase) ID() string {
	return c.Id
}

func (c *CompanyRecordBase) IdDb() uint64 {
	return uint64(c.GetInt("id_db"))
}

func (c *CompanyRecordBase) SetIdDb(id uint64) {
	c.Set("id_db", int(id))
}

func (c *CompanyRecordBase) Name() string {
	return c.GetString("name")
}

func (c *CompanyRecordBase) SetName(name string) {
	c.Set("name", name)
}

func (c *CompanyRecordBase) Checksum() string {
	return c.GetString("checksum")
}

func (c *CompanyRecordBase) SetChecksum(checksum string) {
	c.Set("checksum", checksum)
}

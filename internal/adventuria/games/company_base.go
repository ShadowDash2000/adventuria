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

func (c *CompanyRecordBase) IdDb() string {
	return c.GetString("id_db")
}

func (c *CompanyRecordBase) SetIdDb(id string) {
	c.Set("id_db", id)
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

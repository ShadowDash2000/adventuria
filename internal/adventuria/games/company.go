package games

import (
	"github.com/bestnite/go-igdb/proto"
	"github.com/pocketbase/pocketbase/core"
)

type Company struct {
	core.BaseRecordProxy
}

func (c *Company) IGDBID() int {
	return c.GetInt("igdb_id")
}

func (c *Company) SetIGDBID(id int) {
	c.Set("igdb_id", id)
}

func (c *Company) Name() string {
	return c.GetString("name")
}

func (c *Company) SetName(name string) {
	c.Set("name", name)
}

func (c *Company) Checksum() string {
	return c.GetString("checksum")
}

func (c *Company) SetChecksum(checksum string) {
	c.Set("checksum", checksum)
}

func (c *Company) Data() *proto.Platform {
	var data *proto.Platform
	c.UnmarshalJSONField("data", &data)
	return data
}

func (c *Company) SetData(data *proto.Company) {
	c.Set("data", data)
}

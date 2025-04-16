package games

import (
	"github.com/bestnite/go-igdb/proto"
	"github.com/pocketbase/pocketbase/core"
)

type Platform struct {
	core.BaseRecordProxy
}

func (p *Platform) IGDBID() int {
	return p.GetInt("igdb_id")
}

func (p *Platform) SetIGDBID(id int) {
	p.Set("igdb_id", id)
}

func (p *Platform) Name() string {
	return p.GetString("name")
}

func (p *Platform) SetName(name string) {
	p.Set("name", name)
}

func (p *Platform) Checksum() string {
	return p.GetString("checksum")
}

func (p *Platform) SetChecksum(checksum string) {
	p.Set("checksum", checksum)
}

func (p *Platform) Data() *proto.Platform {
	var data *proto.Platform
	p.UnmarshalJSONField("data", &data)
	return data
}

func (p *Platform) SetData(data *proto.Platform) {
	p.Set("data", data)
}

package steam

import "github.com/pocketbase/pocketbase/core"

type SteamSpyRecord struct {
	core.BaseRecordProxy
}

func NewSteamSpyRecordFromRecord(record *core.Record) *SteamSpyRecord {
	r := &SteamSpyRecord{}
	r.SetProxyRecord(record)
	return r
}

func (r *SteamSpyRecord) IdDb() uint {
	return uint(r.GetInt("id_db"))
}

func (r *SteamSpyRecord) SetIdDb(id uint) {
	r.Set("id_db", id)
}

func (r *SteamSpyRecord) Name() string {
	return r.GetString("name")
}

func (r *SteamSpyRecord) SetName(name string) {
	r.Set("name", name)
}

func (r *SteamSpyRecord) Price() uint {
	return uint(r.GetInt("price"))
}

func (r *SteamSpyRecord) SetPrice(price uint) {
	r.Set("price", price)
}

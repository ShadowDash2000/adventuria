package cheapshark

import "github.com/pocketbase/pocketbase/core"

type CheapSharkResponse struct {
	SteamAppID  uint    `json:"steamAppID"`
	Title       string  `json:"title"`
	NormalPrice float64 `json:"normalPrice"`
}

type CheapSharkRecord struct {
	core.BaseRecordProxy
}

func NewCheapsharkRecordFromRecord(record *core.Record) *CheapSharkRecord {
	r := &CheapSharkRecord{}
	r.SetProxyRecord(record)
	return r
}

func (r *CheapSharkRecord) IdDb() uint {
	return uint(r.GetInt("id_db"))
}

func (r *CheapSharkRecord) SetIdDb(id uint) {
	r.Set("id_db", id)
}

func (r *CheapSharkRecord) Name() string {
	return r.GetString("name")
}

func (r *CheapSharkRecord) SetName(name string) {
	r.Set("name", name)
}

func (r *CheapSharkRecord) Price() float64 {
	return r.GetFloat("price")
}

func (r *CheapSharkRecord) SetPrice(price float64) {
	r.Set("price", price)
}

package hltb

import "github.com/pocketbase/pocketbase/core"

type HowLongToBeatRecord struct {
	core.BaseRecordProxy
}

func NewHowLongToBeatRecordFromRecord(record *core.Record) *HowLongToBeatRecord {
	r := &HowLongToBeatRecord{}
	r.SetProxyRecord(record)
	return r
}

func (r *HowLongToBeatRecord) IdDb() int {
	return r.GetInt("id_db")
}

func (r *HowLongToBeatRecord) SetIdDb(id int) {
	r.Set("id_db", id)
}

func (r *HowLongToBeatRecord) Name() string {
	return r.GetString("name")
}

func (r *HowLongToBeatRecord) SetName(name string) {
	r.Set("name", name)
}

func (r *HowLongToBeatRecord) Campaign() float64 {
	return r.GetFloat("campaign")
}

func (r *HowLongToBeatRecord) SetCampaign(campaign float64) {
	r.Set("campaign", campaign)
}

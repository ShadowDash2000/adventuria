package model

import (
	"errors"
)

type HowLongToBeatData struct {
	Id       string
	IdDb     int
	Name     string
	Year     int
	Campaign float64
}

type HowLongToBeat struct {
	data  HowLongToBeatData
	isNew bool
}

type HowLongToBeatCreate struct {
	IdDb     int
	Name     string
	Year     int
	Campaign float64
}

func NewHowLongToBeat(data HowLongToBeatCreate) (*HowLongToBeat, error) {
	if data.IdDb == 0 {
		return nil, errors.New("howlongtobeat: id_db must be non-zero")
	}
	if data.Name == "" {
		return nil, errors.New("howlongtobeat: name is empty")
	}

	return &HowLongToBeat{
		data: HowLongToBeatData{
			IdDb:     data.IdDb,
			Name:     data.Name,
			Year:     data.Year,
			Campaign: data.Campaign,
		},
		isNew: true,
	}, nil
}

func RestoreHowLongToBeat(data HowLongToBeatData) *HowLongToBeat {
	return &HowLongToBeat{
		data:  data,
		isNew: false,
	}
}

func (h *HowLongToBeat) IsNew() bool {
	return h.isNew
}

func (h *HowLongToBeat) ID() string {
	return h.data.Id
}

func (h *HowLongToBeat) IdDb() int {
	return h.data.IdDb
}

func (h *HowLongToBeat) Name() string {
	return h.data.Name
}

func (h *HowLongToBeat) Year() int {
	return h.data.Year
}

func (h *HowLongToBeat) Campaign() float64 {
	return h.data.Campaign
}

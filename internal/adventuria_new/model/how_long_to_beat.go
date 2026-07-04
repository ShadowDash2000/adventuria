package model

import (
	"errors"

	"github.com/google/uuid"
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

func NewHowLongToBeat(id uuid.UUID, data HowLongToBeatCreate) (*HowLongToBeat, error) {
	if id == uuid.Nil {
		return nil, errors.New("howlongtobeat: id cannot be nil")
	}
	if data.IdDb == 0 {
		return nil, errors.New("howlongtobeat: id_db must be non-zero")
	}
	if data.Name == "" {
		return nil, errors.New("howlongtobeat: name is empty")
	}

	return &HowLongToBeat{
		data: HowLongToBeatData{
			Id:       id.String(),
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

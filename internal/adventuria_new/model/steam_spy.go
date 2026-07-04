package model

import (
	"errors"

	"github.com/google/uuid"
)

type SteamSpyData struct {
	Id    string
	IdDb  int
	Name  string
	Price int
}

type SteamSpy struct {
	data  SteamSpyData
	isNew bool
}

type SteamSpyCreate struct {
	IdDb  int
	Name  string
	Price int
}

func NewSteamSpy(id uuid.UUID, data SteamSpyCreate) (*SteamSpy, error) {
	if id == uuid.Nil {
		return nil, errors.New("steam_spy: id cannot be nil")
	}
	if data.IdDb == 0 {
		return nil, errors.New("steam_spy: id_db must be non-zero")
	}
	if data.Name == "" {
		return nil, errors.New("steam_spy: name is empty")
	}

	return &SteamSpy{
		data: SteamSpyData{
			Id:    id.String(),
			IdDb:  data.IdDb,
			Name:  data.Name,
			Price: data.Price,
		},
		isNew: true,
	}, nil
}

func RestoreSteamSpy(data SteamSpyData) *SteamSpy {
	return &SteamSpy{
		data:  data,
		isNew: false,
	}
}

func (s *SteamSpy) IsNew() bool {
	return s.isNew
}

func (s *SteamSpy) ID() string {
	return s.data.Id
}

func (s *SteamSpy) IdDb() int {
	return s.data.IdDb
}

func (s *SteamSpy) Name() string {
	return s.data.Name
}

func (s *SteamSpy) Price() int {
	return s.data.Price
}

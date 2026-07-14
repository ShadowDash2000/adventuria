package model

import (
	"errors"
)

type GameTypeData struct {
	Id       string
	IdDb     string
	Name     string
	Checksum string
}

type GameType struct {
	data  GameTypeData
	isNew bool
}

type GameTypeCreate struct {
	IdDb     string
	Name     string
	Checksum string
}

func NewGameType(data GameTypeCreate) (*GameType, error) {
	if data.IdDb == "" {
		return nil, errors.New("game_type: idDb is empty")
	}
	if data.Name == "" {
		return nil, errors.New("game_type: name is empty")
	}
	if data.Checksum == "" {
		return nil, errors.New("game_type: checksum is empty")
	}

	return &GameType{
		data: GameTypeData{
			IdDb:     data.IdDb,
			Name:     data.Name,
			Checksum: data.Checksum,
		},
		isNew: true,
	}, nil
}

func RestoreGameType(data GameTypeData) *GameType {
	return &GameType{
		data:  data,
		isNew: false,
	}
}

func (p *GameType) IsNew() bool {
	return p.isNew
}

func (p *GameType) ID() string {
	return p.data.Id
}

func (p *GameType) IdDb() string {
	return p.data.IdDb
}

func (p *GameType) Name() string {
	return p.data.Name
}

func (p *GameType) SetName(name string) {
	p.data.Name = name
}

func (p *GameType) Checksum() string {
	return p.data.Checksum
}

func (p *GameType) SetChecksum(checksum string) {
	p.data.Checksum = checksum
}

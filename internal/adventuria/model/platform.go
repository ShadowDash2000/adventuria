package model

import (
	"errors"
)

type PlatformData struct {
	Id       string
	IdDb     string
	Name     string
	Checksum string
}

type Platform struct {
	data  PlatformData
	isNew bool
}

type PlatformCreate struct {
	IdDb     string
	Name     string
	Checksum string
}

func NewPlatform(data PlatformCreate) (*Platform, error) {
	if data.IdDb == "" {
		return nil, errors.New("platform: idDb is empty")
	}
	if data.Name == "" {
		return nil, errors.New("platform: name is empty")
	}
	if data.Checksum == "" {
		return nil, errors.New("platform: checksum is empty")
	}

	return &Platform{
		data: PlatformData{
			IdDb:     data.IdDb,
			Name:     data.Name,
			Checksum: data.Checksum,
		},
		isNew: true,
	}, nil
}

func RestorePlatform(data PlatformData) *Platform {
	return &Platform{
		data:  data,
		isNew: false,
	}
}

func (p *Platform) IsNew() bool {
	return p.isNew
}

func (p *Platform) ID() string {
	return p.data.Id
}

func (p *Platform) IdDb() string {
	return p.data.IdDb
}

func (p *Platform) Name() string {
	return p.data.Name
}

func (p *Platform) SetName(name string) {
	p.data.Name = name
}

func (p *Platform) Checksum() string {
	return p.data.Checksum
}

func (p *Platform) SetChecksum(checksum string) {
	p.data.Checksum = checksum
}

func (p *Platform) SetIdDb(idDb string) {
	p.data.IdDb = idDb
}

package model

import (
	"errors"

	"github.com/google/uuid"
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

func NewPlatform(id uuid.UUID, data PlatformCreate) (*Platform, error) {
	if id == uuid.Nil {
		return nil, errors.New("platform: id cannot be nil")
	}
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
			Id:       id.String(),
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

func (p *Platform) Checksum() string {
	return p.data.Checksum
}

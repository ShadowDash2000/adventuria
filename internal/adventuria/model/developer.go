package model

import (
	"errors"

	"github.com/google/uuid"
)

type DeveloperData struct {
	Id       string
	IdDb     string
	Name     string
	Checksum string
}

type Developer struct {
	data  DeveloperData
	isNew bool
}

type DeveloperCreate struct {
	IdDb     string
	Name     string
	Checksum string
}

func NewDeveloper(id uuid.UUID, data DeveloperCreate) (*Developer, error) {
	if id == uuid.Nil {
		return nil, errors.New("developer: id cannot be nil")
	}
	if data.IdDb == "" {
		return nil, errors.New("developer: idDb is empty")
	}
	if data.Name == "" {
		return nil, errors.New("developer: name is empty")
	}
	if data.Checksum == "" {
		return nil, errors.New("developer: checksum is empty")
	}

	return &Developer{
		data: DeveloperData{
			Id:       id.String(),
			IdDb:     data.IdDb,
			Name:     data.Name,
			Checksum: data.Checksum,
		},
		isNew: true,
	}, nil
}

func RestoreDeveloper(data DeveloperData) *Developer {
	return &Developer{
		data:  data,
		isNew: false,
	}
}

func (d *Developer) IsNew() bool {
	return d.isNew
}

func (d *Developer) ID() string {
	return d.data.Id
}

func (d *Developer) IdDb() string {
	return d.data.IdDb
}

func (d *Developer) Name() string {
	return d.data.Name
}

func (d *Developer) Checksum() string {
	return d.data.Checksum
}

package model

import (
	"errors"

	"github.com/google/uuid"
)

type PublisherData struct {
	Id       string
	IdDb     string
	Name     string
	Checksum string
}

type Publisher struct {
	data  PublisherData
	isNew bool
}

type PublisherCreate struct {
	IdDb     string
	Name     string
	Checksum string
}

func NewPublisher(id uuid.UUID, data PublisherCreate) (*Publisher, error) {
	if id == uuid.Nil {
		return nil, errors.New("publisher: id cannot be nil")
	}
	if data.IdDb == "" {
		return nil, errors.New("publisher: idDb is empty")
	}
	if data.Name == "" {
		return nil, errors.New("publisher: name is empty")
	}
	if data.Checksum == "" {
		return nil, errors.New("publisher: checksum is empty")
	}

	return &Publisher{
		data: PublisherData{
			Id:       id.String(),
			IdDb:     data.IdDb,
			Name:     data.Name,
			Checksum: data.Checksum,
		},
		isNew: true,
	}, nil
}

func RestorePublisher(data PublisherData) *Publisher {
	return &Publisher{
		data:  data,
		isNew: false,
	}
}

func (p *Publisher) IsNew() bool {
	return p.isNew
}

func (p *Publisher) ID() string {
	return p.data.Id
}

func (p *Publisher) IdDb() string {
	return p.data.IdDb
}

func (p *Publisher) Name() string {
	return p.data.Name
}

func (p *Publisher) Checksum() string {
	return p.data.Checksum
}

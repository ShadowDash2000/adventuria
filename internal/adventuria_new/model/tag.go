package model

import (
	"errors"

	"github.com/google/uuid"
)

type TagData struct {
	Id       string
	IdDb     string
	Name     string
	Checksum string
}

type Tag struct {
	data  TagData
	isNew bool
}

type TagCreate struct {
	IdDb     string
	Name     string
	Checksum string
}

func NewTag(id uuid.UUID, data TagCreate) (*Tag, error) {
	if id == uuid.Nil {
		return nil, errors.New("tag: id cannot be nil")
	}
	if data.IdDb == "" {
		return nil, errors.New("tag: idDb is empty")
	}
	if data.Name == "" {
		return nil, errors.New("tag: name is empty")
	}
	if data.Checksum == "" {
		return nil, errors.New("tag: checksum is empty")
	}

	return &Tag{
		data: TagData{
			Id:       id.String(),
			IdDb:     data.IdDb,
			Name:     data.Name,
			Checksum: data.Checksum,
		},
		isNew: true,
	}, nil
}

func RestoreTag(data TagData) *Tag {
	return &Tag{
		data:  data,
		isNew: false,
	}
}

func (t *Tag) IsNew() bool {
	return t.isNew
}

func (t *Tag) ID() string {
	return t.data.Id
}

func (t *Tag) IdDb() string {
	return t.data.IdDb
}

func (t *Tag) Name() string {
	return t.data.Name
}

func (t *Tag) Checksum() string {
	return t.data.Checksum
}

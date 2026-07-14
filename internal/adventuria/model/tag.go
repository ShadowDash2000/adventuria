package model

import (
	"errors"
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

func NewTag(data TagCreate) (*Tag, error) {
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

func (t *Tag) SetName(name string) {
	t.data.Name = name
}

func (t *Tag) Checksum() string {
	return t.data.Checksum
}

func (t *Tag) SetChecksum(checksum string) {
	t.data.Checksum = checksum
}

package model

import (
	"errors"
)

type ThemeData struct {
	Id       string
	IdDb     string
	Name     string
	Checksum string
}

type Theme struct {
	data  ThemeData
	isNew bool
}

type ThemeCreate struct {
	IdDb     string
	Name     string
	Checksum string
}

func NewTheme(data ThemeCreate) (*Theme, error) {
	if data.IdDb == "" {
		return nil, errors.New("theme: idDb is empty")
	}
	if data.Name == "" {
		return nil, errors.New("theme: name is empty")
	}
	if data.Checksum == "" {
		return nil, errors.New("theme: checksum is empty")
	}

	return &Theme{
		data: ThemeData{
			IdDb:     data.IdDb,
			Name:     data.Name,
			Checksum: data.Checksum,
		},
		isNew: true,
	}, nil
}

func RestoreTheme(data ThemeData) *Theme {
	return &Theme{
		data:  data,
		isNew: false,
	}
}

func (t *Theme) IsNew() bool {
	return t.isNew
}

func (t *Theme) ID() string {
	return t.data.Id
}

func (t *Theme) IdDb() string {
	return t.data.IdDb
}

func (t *Theme) Name() string {
	return t.data.Name
}

func (t *Theme) SetName(name string) {
	t.data.Name = name
}

func (t *Theme) Checksum() string {
	return t.data.Checksum
}

func (t *Theme) SetChecksum(checksum string) {
	t.data.Checksum = checksum
}

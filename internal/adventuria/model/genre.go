package model

import (
	"errors"
)

type GenreData struct {
	Id       string
	IdDb     string
	Name     string
	Checksum string
}

type Genre struct {
	data  GenreData
	isNew bool
}

type GenreCreate struct {
	IdDb     string
	Name     string
	Checksum string
}

func NewGenre(data GenreCreate) (*Genre, error) {
	if data.IdDb == "" {
		return nil, errors.New("genre: idDb is empty")
	}
	if data.Name == "" {
		return nil, errors.New("genre: name is empty")
	}
	if data.Checksum == "" {
		return nil, errors.New("genre: checksum is empty")
	}

	return &Genre{
		data: GenreData{
			IdDb:     data.IdDb,
			Name:     data.Name,
			Checksum: data.Checksum,
		},
		isNew: true,
	}, nil
}

func RestoreGenre(data GenreData) *Genre {
	return &Genre{
		data:  data,
		isNew: false,
	}
}

func (g *Genre) IsNew() bool {
	return g.isNew
}

func (g *Genre) ID() string {
	return g.data.Id
}

func (g *Genre) IdDb() string {
	return g.data.IdDb
}

func (g *Genre) Name() string {
	return g.data.Name
}

func (g *Genre) SetName(name string) {
	g.data.Name = name
}

func (g *Genre) Checksum() string {
	return g.data.Checksum
}

func (g *Genre) SetChecksum(checksum string) {
	g.data.Checksum = checksum
}

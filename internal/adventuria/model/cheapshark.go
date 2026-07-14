package model

import (
	"errors"
)

type CheapSharkData struct {
	Id    string
	IdDb  int
	Name  string
	Price float64
}

type CheapShark struct {
	data  CheapSharkData
	isNew bool
}

type CheapSharkCreate struct {
	IdDb  int
	Name  string
	Price float64
}

func NewCheapShark(data CheapSharkCreate) (*CheapShark, error) {
	if data.IdDb == 0 {
		return nil, errors.New("cheapshark: id_db must be non-zero")
	}
	if data.Name == "" {
		return nil, errors.New("cheapshark: name is empty")
	}
	if data.Price <= 0 {
		return nil, errors.New("cheapshark: price must be greater than zero")
	}

	return &CheapShark{
		data: CheapSharkData{
			IdDb:  data.IdDb,
			Name:  data.Name,
			Price: data.Price,
		},
		isNew: true,
	}, nil
}

func RestoreCheapShark(data CheapSharkData) *CheapShark {
	return &CheapShark{
		data:  data,
		isNew: false,
	}
}

func (c *CheapShark) IsNew() bool {
	return c.isNew
}

func (c *CheapShark) ID() string {
	return c.data.Id
}

func (c *CheapShark) IdDb() int {
	return c.data.IdDb
}

func (c *CheapShark) Name() string {
	return c.data.Name
}

func (c *CheapShark) Price() float64 {
	return c.data.Price
}

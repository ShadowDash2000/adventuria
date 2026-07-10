package model

import (
	"errors"

	"github.com/google/uuid"
)

type CompanyData struct {
	Id       string
	IdDb     string
	Name     string
	Checksum string
}

type Company struct {
	data  CompanyData
	isNew bool
}

type CompanyCreate struct {
	IdDb     string
	Name     string
	Checksum string
}

func NewCompany(id uuid.UUID, data CompanyCreate) (*Company, error) {
	if id == uuid.Nil {
		return nil, errors.New("company: id cannot be nil")
	}
	if data.IdDb == "" {
		return nil, errors.New("company: idDb is empty")
	}
	if data.Name == "" {
		return nil, errors.New("company: name is empty")
	}
	if data.Checksum == "" {
		return nil, errors.New("company: checksum is empty")
	}

	return &Company{
		data: CompanyData{
			Id:       id.String(),
			IdDb:     data.IdDb,
			Name:     data.Name,
			Checksum: data.Checksum,
		},
		isNew: true,
	}, nil
}

func RestoreCompany(data CompanyData) *Company {
	return &Company{
		data:  data,
		isNew: false,
	}
}

func (c *Company) IsNew() bool {
	return c.isNew
}

func (c *Company) ID() string {
	return c.data.Id
}

func (c *Company) IdDb() string {
	return c.data.IdDb
}

func (c *Company) Name() string {
	return c.data.Name
}

func (c *Company) SetName(name string) {
	c.data.Name = name
}

func (c *Company) Checksum() string {
	return c.data.Checksum
}

func (c *Company) SetChecksum(checksum string) {
	c.data.Checksum = checksum
}

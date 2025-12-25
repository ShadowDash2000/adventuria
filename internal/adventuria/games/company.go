package games

import "github.com/pocketbase/pocketbase/core"

type Company struct {
	IdDb     string
	Name     string
	Checksum string
}

type CompanyRecord interface {
	core.RecordProxy
	UpdatableRecord

	ID() string
	IdDb() string
	SetIdDb(string)
	Name() string
	SetName(string)
}

package games

import (
	"github.com/pocketbase/pocketbase/core"
)

type Company struct {
	IdDb     uint64
	Name     string
	Checksum string
}

type CompanyRecord interface {
	core.RecordProxy

	ID() string
	IdDb() uint64
	SetIdDb(int)
	Name() string
	SetName(string)

	Checksum() string
	SetChecksum(string)
}

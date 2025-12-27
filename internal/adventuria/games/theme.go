package games

import "github.com/pocketbase/pocketbase/core"

type Theme struct {
	IdDb     string
	Name     string
	Checksum string
}

type ThemeRecord interface {
	core.RecordProxy
	UpdatableRecord

	ID() string
	SetIdDb(string)
	Name() string
	SetName(string)
}

package igdb

import (
	"time"
)

type TableReference struct {
	Ids        []uint64
	TableName  string
	PrimaryKey string
	SearchKey  string
}

type TableReferenceSingle struct {
	Id         uint64
	TableName  string
	PrimaryKey string
	SearchKey  string
}

type Game struct {
	Id          string
	Name        string
	Slug        string
	ReleaseDate time.Time
	Platforms   []uint64
	Developers  []uint64
	Publishers  []uint64
	Genres      []uint64
	Tags        []uint64
	Themes      []uint64
	GameType    uint64
	SteamAppId  uint64
	Cover       string
	Checksum    string
}

type Company struct {
	Id       string
	Name     string
	Checksum string
}

type Tag struct {
	Id       string
	Name     string
	Checksum string
}

type Theme struct {
	Id       string
	Name     string
	Checksum string
}

type Platform struct {
	Id       string
	Name     string
	Checksum string
}

type Genre struct {
	Id       string
	Name     string
	Checksum string
}

type GameType struct {
	Id       string
	Name     string
	Checksum string
}

type ParseGamesMessage struct {
	Games []*Game
	Err   error
}

type ParsePlatformsMessage struct {
	Platforms []*Platform
	Err       error
}

type ParseGenresMessage struct {
	Genres []*Genre
	Err    error
}

type ParseGameTypesMessage struct {
	GameTypes []*GameType
	Err       error
}

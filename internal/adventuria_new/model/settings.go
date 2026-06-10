package model

import (
	"errors"

	"github.com/google/uuid"
)

type SettingsData struct {
	Id                string
	EventEnded        bool
	CurrentSeason     string
	CurrentWeek       int
	BlockAllActions   bool
	MaxInventorySlots int
	PointsForDrop     int
	DropsToJail       int

	IgdbGamesParsed         int
	DisableIgdbParser       bool
	DisableSteamParser      bool
	DisableCheapsharkParser bool
	DisableHltbParser       bool
	DisableRefreshHltbTime  bool
	KillParser              bool
	IgdbForceUpdateGames    bool
}

type Settings struct {
	data  SettingsData
	isNew bool
}

type SettingsCreate struct {
	CurrentSeason string
	DropsToJail   int
}

func NewSettings(id uuid.UUID, data SettingsCreate) (*Settings, error) {
	if id == uuid.Nil {
		return nil, errors.New("settings: id cannot be nil")
	}
	if data.CurrentSeason == "" {
		return nil, errors.New("settings: current season is empty")
	}
	if data.DropsToJail < 0 {
		return nil, errors.New("settings: drops to jail cannot be negative")
	}

	return &Settings{
		data: SettingsData{
			Id:            id.String(),
			EventEnded:    false,
			CurrentSeason: data.CurrentSeason,
			DropsToJail:   data.DropsToJail,
		},
		isNew: true,
	}, nil
}

func RestoreSettings(data SettingsData) *Settings {
	return &Settings{
		data:  data,
		isNew: false,
	}
}

func (s *Settings) IsNew() bool {
	return s.isNew
}

func (s *Settings) Id() string {
	return s.data.Id
}

func (s *Settings) EventEnded() bool {
	return s.data.EventEnded
}

func (s *Settings) CurrentSeason() string {
	return s.data.CurrentSeason
}

func (s *Settings) CurrentWeek() int {
	return s.data.CurrentWeek
}

func (s *Settings) BlockAllActions() bool {
	return s.data.BlockAllActions
}

func (s *Settings) MaxInventorySlots() int {
	return s.data.MaxInventorySlots
}

func (s *Settings) PointsForDrop() int {
	return s.data.PointsForDrop
}

func (s *Settings) DropsToJail() int {
	return s.data.DropsToJail
}

func (s *Settings) IgdbGamesParsed() int {
	return s.data.IgdbGamesParsed
}

func (s *Settings) DisableIgdbParser() bool {
	return s.data.DisableIgdbParser
}

func (s *Settings) DisableSteamParser() bool {
	return s.data.DisableSteamParser
}

func (s *Settings) DisableCheapsharkParser() bool {
	return s.data.DisableCheapsharkParser
}

func (s *Settings) DisableHltbParser() bool {
	return s.data.DisableHltbParser
}

func (s *Settings) DisableRefreshHltbTime() bool {
	return s.data.DisableRefreshHltbTime
}

func (s *Settings) KillParser() bool {
	return s.data.KillParser
}

func (s *Settings) IgdbForceUpdateGames() bool {
	return s.data.IgdbForceUpdateGames
}

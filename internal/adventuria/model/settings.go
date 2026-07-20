package model

import (
	"errors"
	"time"
)

type SettingsData struct {
	Id                string
	EventEnded        bool
	CurrentSeason     string
	CurrentWeek       int
	BlockAllActions   bool
	EnergyDefault     int
	MaxInventorySlots int
	PointsForDrop     int
	DropsToJail       int

	IgdbFilter IgdbFilter

	IgdbGamesParsed         uint
	DisableIgdbParser       bool
	DisableIgdbGamesParser  bool
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

func NewSettings(data SettingsCreate) (*Settings, error) {
	if data.CurrentSeason == "" {
		return nil, errors.New("settings: current season is empty")
	}
	if data.DropsToJail < 0 {
		return nil, errors.New("settings: drops to jail cannot be negative")
	}

	return &Settings{
		data: SettingsData{
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

func (s *Settings) SetEventEnded(b bool) {
	s.data.EventEnded = b
}

func (s *Settings) CurrentSeason() string {
	return s.data.CurrentSeason
}

func (s *Settings) SetCurrentSeason(season string) {
	s.data.CurrentSeason = season
}

func (s *Settings) CurrentWeek() int {
	return s.data.CurrentWeek
}

func (s *Settings) SetCurrentWeek(week int) {
	s.data.CurrentWeek = week
}

func (s *Settings) BlockAllActions() bool {
	return s.data.BlockAllActions
}

func (s *Settings) SetBlockAllActions(b bool) {
	s.data.BlockAllActions = b
}

func (s *Settings) EnergyDefault() int {
	return s.data.EnergyDefault
}

func (s *Settings) SetEnergyDefault(energy int) {
	s.data.EnergyDefault = energy
}

func (s *Settings) MaxInventorySlots() int {
	return s.data.MaxInventorySlots
}

func (s *Settings) SetMaxInventorySlots(slots int) {
	s.data.MaxInventorySlots = slots
}

func (s *Settings) PointsForDrop() int {
	return s.data.PointsForDrop
}

func (s *Settings) SetPointsForDrop(points int) {
	s.data.PointsForDrop = points
}

func (s *Settings) DropsToJail() int {
	return s.data.DropsToJail
}

func (s *Settings) SetDropsToJail(drops int) {
	s.data.DropsToJail = drops
}

func (s *Settings) IgdbFilter() IgdbFilter {
	return s.data.IgdbFilter
}

func (s *Settings) IgdbFilterGameTypes() []string {
	return s.data.IgdbFilter.GameTypes
}

func (s *Settings) SetIgdbFilterGameTypes(types []string) {
	s.data.IgdbFilter.GameTypes = types
}

func (s *Settings) IgdbFilterPlatforms() []string {
	return s.data.IgdbFilter.Platforms
}

func (s *Settings) SetIgdbFilterPlatforms(platforms []string) {
	s.data.IgdbFilter.Platforms = platforms
}

func (s *Settings) IgdbFilterFirstReleaseDateMin() time.Time {
	return s.data.IgdbFilter.ReleaseDateMin
}

func (s *Settings) SetIgdbFilterFirstReleaseDateMin(t time.Time) {
	s.data.IgdbFilter.ReleaseDateMin = t
}

func (s *Settings) IgdbFilterFirstReleaseDateMax() time.Time {
	return s.data.IgdbFilter.ReleaseDateMax
}

func (s *Settings) SetIgdbFilterFirstReleaseDateMax(t time.Time) {
	s.data.IgdbFilter.ReleaseDateMax = t
}

func (s *Settings) IgdbGamesParsed() uint {
	return s.data.IgdbGamesParsed
}

func (s *Settings) SetIgdbGamesParsed(n uint) {
	s.data.IgdbGamesParsed = n
}

func (s *Settings) DisableIgdbParser() bool {
	return s.data.DisableIgdbParser
}

func (s *Settings) SetDisableIgdbParser(b bool) {
	s.data.DisableIgdbParser = b
}

func (s *Settings) DisableIgdbGamesParser() bool {
	return s.data.DisableIgdbGamesParser
}

func (s *Settings) SetDisableIgdbGamesParser(b bool) {
	s.data.DisableIgdbGamesParser = b
}

func (s *Settings) DisableSteamParser() bool {
	return s.data.DisableSteamParser
}

func (s *Settings) SetDisableSteamParser(b bool) {
	s.data.DisableSteamParser = b
}

func (s *Settings) DisableCheapsharkParser() bool {
	return s.data.DisableCheapsharkParser
}

func (s *Settings) SetDisableCheapsharkParser(b bool) {
	s.data.DisableCheapsharkParser = b
}

func (s *Settings) DisableHltbParser() bool {
	return s.data.DisableHltbParser
}

func (s *Settings) SetDisableHltbParser(b bool) {
	s.data.DisableHltbParser = b
}

func (s *Settings) DisableRefreshHltbTime() bool {
	return s.data.DisableRefreshHltbTime
}

func (s *Settings) SetDisableRefreshHltbTime(b bool) {
	s.data.DisableRefreshHltbTime = b
}

func (s *Settings) KillParser() bool {
	return s.data.KillParser
}

func (s *Settings) IgdbForceUpdateGames() bool {
	return s.data.IgdbForceUpdateGames
}

func (s *Settings) SetIgdbForceUpdateGames(b bool) {
	s.data.IgdbForceUpdateGames = b
}

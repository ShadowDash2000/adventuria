package model

import (
	"errors"
)

type PlayerStatsData struct {
	Id              string
	Player          string
	Season          string
	ActivitiesStats ActivitiesStats
	CellsPassed     int
	Drops           int
	Rerolls         int
	WasInJail       int
	ItemsUsed       int
	DiceRolls       int
	MaxDiceRoll     int
	WheelsRolled    int
}

type ActivitiesStats struct {
	GamesCompleted   int
	MoviesCompleted  int
	GymsCompleted    int
	KaraokeCompleted int
}

type PlayerStats struct {
	data  PlayerStatsData
	isNew bool
}

type PlayerStatsCreate struct {
	Player string
	Season string
}

func NewPlayerStats(data PlayerStatsCreate) (*PlayerStats, error) {
	if data.Player == "" {
		return nil, errors.New("player is empty")
	}
	if data.Season == "" {
		return nil, errors.New("season is empty")
	}

	return &PlayerStats{
		data: PlayerStatsData{
			Player: data.Player,
			Season: data.Season,
		},
		isNew: true,
	}, nil
}

func RestorePlayerStats(data PlayerStatsData) *PlayerStats {
	return &PlayerStats{
		data:  data,
		isNew: false,
	}
}

func (p *PlayerStats) IsNew() bool {
	return p.isNew
}

func (p *PlayerStats) ID() string {
	return p.data.Id
}

func (p *PlayerStats) Player() string {
	return p.data.Player
}

func (p *PlayerStats) Season() string {
	return p.data.Season
}

func (p *PlayerStats) ActivitiesStats() ActivitiesStats {
	return p.data.ActivitiesStats
}

func (p *PlayerStats) GamesCompletedChange(amount int) {
	p.data.ActivitiesStats.GamesCompleted += amount
}

func (p *PlayerStats) MoviesCompletedChange(amount int) {
	p.data.ActivitiesStats.MoviesCompleted += amount
}

func (p *PlayerStats) GymsCompletedChange(amount int) {
	p.data.ActivitiesStats.GymsCompleted += amount
}

func (p *PlayerStats) KaraokeCompletedChange(amount int) {
	p.data.ActivitiesStats.KaraokeCompleted += amount
}

func (p *PlayerStats) CellsPassed() int {
	return p.data.CellsPassed
}

func (p *PlayerStats) CellsPassedChange(amount int) {
	p.data.CellsPassed += amount
}

func (p *PlayerStats) Drops() int {
	return p.data.Drops
}

func (p *PlayerStats) DropsChange(amount int) {
	p.data.Drops += amount
}

func (p *PlayerStats) Rerolls() int {
	return p.data.Rerolls
}

func (p *PlayerStats) RerollsChange(amount int) {
	p.data.Rerolls += amount
}

func (p *PlayerStats) WasInJail() int {
	return p.data.WasInJail
}

func (p *PlayerStats) WasInJailChange(amount int) {
	p.data.WasInJail += amount
}

func (p *PlayerStats) ItemsUsed() int {
	return p.data.ItemsUsed
}

func (p *PlayerStats) ItemsUsedChange(amount int) {
	p.data.ItemsUsed += amount
}

func (p *PlayerStats) DiceRolls() int {
	return p.data.DiceRolls
}

func (p *PlayerStats) DiceRollsChange(amount int) {
	p.data.DiceRolls += amount
}

func (p *PlayerStats) MaxDiceRoll() int {
	return p.data.MaxDiceRoll
}

func (p *PlayerStats) MaxDiceRollChange(amount int) {
	p.data.MaxDiceRoll += amount
}

func (p *PlayerStats) WheelsRolled() int {
	return p.data.WheelsRolled
}

func (p *PlayerStats) WheelsRolledChange(amount int) {
	p.data.WheelsRolled += amount
}

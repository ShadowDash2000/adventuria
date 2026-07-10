package model

import (
	"errors"

	"github.com/google/uuid"
)

type PlayerProgressData struct {
	Id                string
	Player            string
	Season            string
	CurrentWorld      string
	Points            int
	Balance           int
	CellsPassed       int
	IsInJail          bool
	DropsInARow       int
	ItemWheelsCount   int
	MaxInventorySlots int
}

type PlayerProgress struct {
	data  PlayerProgressData
	isNew bool
}

type PlayerProgressCreate struct {
	Player            string
	Season            string
	CurrentWorld      string
	MaxInventorySlots int
}

func NewPlayerProgress(id uuid.UUID, data PlayerProgressCreate) (*PlayerProgress, error) {
	if id == uuid.Nil {
		return nil, errors.New("player_progress: id cannot be nil")
	}
	if data.Player == "" {
		return nil, errors.New("player_progress: player is empty")
	}
	if data.Season == "" {
		return nil, errors.New("player_progress: season is empty")
	}
	if data.CurrentWorld == "" {
		return nil, errors.New("player_progress: current world is empty")
	}
	if data.MaxInventorySlots < 0 {
		return nil, errors.New("player_progress: max inventory slots cannot be negative")
	}

	return &PlayerProgress{
		data: PlayerProgressData{
			Id:                id.String(),
			Player:            data.Player,
			Season:            data.Season,
			CurrentWorld:      data.CurrentWorld,
			MaxInventorySlots: data.MaxInventorySlots,
		},
		isNew: true,
	}, nil
}

func RestorePlayerProgress(data PlayerProgressData) *PlayerProgress {
	return &PlayerProgress{
		data:  data,
		isNew: false,
	}
}

func (p *PlayerProgress) IsNew() bool {
	return p.isNew
}

func (p *PlayerProgress) ID() string {
	return p.data.Id
}

func (p *PlayerProgress) Player() string {
	return p.data.Player
}

func (p *PlayerProgress) Season() string {
	return p.data.Season
}

func (p *PlayerProgress) CurrentWorld() string {
	return p.data.CurrentWorld
}

func (p *PlayerProgress) SetCurrentWorld(world string) {
	p.data.CurrentWorld = world
}

func (p *PlayerProgress) Points() int {
	return p.data.Points
}

func (p *PlayerProgress) PointsChange(amount int) error {
	if amount == 0 {
		return nil
	}

	p.data.Points += amount

	return nil
}

func (p *PlayerProgress) Balance() int {
	return p.data.Balance
}

func (p *PlayerProgress) BalanceChange(amount int) error {
	if amount == 0 {
		return nil
	}

	p.data.Balance += amount

	return nil
}

func (p *PlayerProgress) CellsPassed() int {
	return p.data.CellsPassed
}

func (p *PlayerProgress) CellsPassedChange(amount int) error {
	if amount == 0 {
		return nil
	}

	p.data.CellsPassed += amount

	return nil
}

func (p *PlayerProgress) SetCellsPassed(count int) {
	p.data.CellsPassed = count
}

func (p *PlayerProgress) IsInJail() bool {
	return p.data.IsInJail
}

func (p *PlayerProgress) SetIsInJail(isInJail bool) {
	p.data.IsInJail = isInJail
}

func (p *PlayerProgress) DropsInARow() int {
	return p.data.DropsInARow
}

func (p *PlayerProgress) DropsInARowChange(amount int) error {
	if amount == 0 {
		return nil
	}

	p.data.DropsInARow += amount

	return nil
}

func (p *PlayerProgress) SetDropsInARow(count int) {
	p.data.DropsInARow = count
}

func (p *PlayerProgress) ItemWheelsCount() int {
	return p.data.ItemWheelsCount
}

func (p *PlayerProgress) ItemWheelsCountChange(amount int) error {
	if amount == 0 {
		return nil
	}

	p.data.ItemWheelsCount += amount

	return nil
}

func (p *PlayerProgress) MaxInventorySlots() int {
	return p.data.MaxInventorySlots
}

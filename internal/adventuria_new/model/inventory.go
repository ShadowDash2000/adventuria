package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type InventoryData struct {
	Id             string
	Activated      time.Time
	Player         string
	Item           string
	IsActive       bool
	AppliedEffects []string
}

type Inventory struct {
	data  InventoryData
	isNew bool
}

type InventoryCreate struct {
	Player    string
	Item      string
	Activated time.Time
	IsActive  bool
}

func NewInventory(id uuid.UUID, data InventoryCreate) (*Inventory, error) {
	if id == uuid.Nil {
		return nil, errors.New("inventory_item: id cannot be nil")
	}
	if data.Player == "" {
		return nil, errors.New("inventory_item: player is empty")
	}
	if data.Item == "" {
		return nil, errors.New("inventory_item: item is empty")
	}
	if data.IsActive && data.Activated.IsZero() {
		return nil, errors.New("inventory_item: activated time cannot be zero when item is active")
	}

	return &Inventory{
		data: InventoryData{
			Id:        id.String(),
			Player:    data.Player,
			Item:      data.Item,
			Activated: data.Activated,
			IsActive:  data.IsActive,
		},
		isNew: true,
	}, nil
}

func RestoreInventory(data InventoryData) *Inventory {
	return &Inventory{
		data:  data,
		isNew: false,
	}
}

func (i *Inventory) IsNew() bool {
	return i.isNew
}

func (i *Inventory) ID() string {
	return i.data.Id
}

func (i *Inventory) Activated() time.Time {
	return i.data.Activated
}

func (i *Inventory) SetActivated(t time.Time) {
	i.data.Activated = t
}

func (i *Inventory) Player() string {
	return i.data.Player
}

func (i *Inventory) Item() string {
	return i.data.Item
}

func (i *Inventory) IsActive() bool {
	return i.data.IsActive
}

func (i *Inventory) SetIsActive(b bool) {
	i.data.IsActive = b
}

func (i *Inventory) AppliedEffects() []string {
	return i.data.AppliedEffects
}

func (i *Inventory) AddAppliedEffects(effects ...string) {
	i.data.AppliedEffects = append(i.data.AppliedEffects, effects...)
}

func (i *Inventory) AppliedEffectsCount() int {
	return len(i.AppliedEffects())
}

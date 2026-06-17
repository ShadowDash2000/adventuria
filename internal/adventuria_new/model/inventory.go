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
	Player   string
	Item     string
	IsActive bool
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

	i := &Inventory{
		data: InventoryData{
			Id:     id.String(),
			Player: data.Player,
			Item:   data.Item,
		},
		isNew: true,
	}

	if data.IsActive {
		i.data.IsActive = true
		i.data.Activated = time.Now()
	}

	return i, nil
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

func (i *Inventory) Activate() error {
	if i.data.IsActive {
		return errors.New("item is already active")
	}
	i.data.IsActive = true
	i.data.Activated = time.Now()
	return nil
}

func (i *Inventory) Activated() time.Time {
	return i.data.Activated
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

func (i *Inventory) AppliedEffects() []string {
	return i.data.AppliedEffects
}

func (i *Inventory) AddAppliedEffects(effects ...string) {
	i.data.AppliedEffects = append(i.data.AppliedEffects, effects...)
}

func (i *Inventory) AppliedEffectsCount() int {
	return len(i.AppliedEffects())
}

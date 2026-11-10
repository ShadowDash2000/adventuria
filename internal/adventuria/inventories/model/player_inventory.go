package model

import (
	"adventuria/internal/adventuria/model"
	"time"
)

type PlayerInventory struct {
	PlayerId           string
	PlayerName         string
	Balance            int
	MaxInventorySlots  int
	UsedInventorySlots int
	Items              []*PlayerInventoryItem
}

type PlayerInventoryItem struct {
	InventoryId      string
	ItemId           string
	Created          time.Time
	Name             string
	Icon             string
	Description      string
	IsActive         bool
	CanDrop          bool
	CanUse           bool
	DropPrice        int
	UnappliedEffects []*PlayerInventoryItemEffect
}

type PlayerInventoryItemEffect struct {
	Id   string
	Type model.EffectType
}

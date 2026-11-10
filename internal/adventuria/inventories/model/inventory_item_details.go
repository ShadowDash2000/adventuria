package model

import (
	"adventuria/internal/adventuria/model"
	"time"
)

type InventoryItemDetails struct {
	Created   time.Time
	IconPath  string
	Inventory *model.Inventory
	Item      *model.Item
	Effects   []*model.EffectInfo
}

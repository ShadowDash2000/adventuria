package adventuria

import "adventuria/pkg/cache"

type ItemRecord interface {
	ID() string
	Name() string
	Icon() string
	IsUsingSlot() bool
	IsActiveByDefault() bool
	CanDrop() bool
	IsRollable() bool
	Order() int
	Type() string
	Price() int
}

type ItemType string

var (
	ItemTypeBuff    ItemType = "buff"
	ItemTypeDebuff  ItemType = "debuff"
	ItemTypeNeutral ItemType = "neutral"
)

var ItemTypes = map[ItemType]bool{
	ItemTypeBuff:    true,
	ItemTypeDebuff:  true,
	ItemTypeNeutral: true,
}

type Item interface {
	ItemRecord
	cache.Closable

	IDInventory() string
	IsActive() bool
	EffectsCount() int
	AppliedEffectsCount() int
	Use() (OnUseSuccess, OnUseFail, error)
	Drop() error
}

type OnUseSuccess func() error
type OnUseFail func()

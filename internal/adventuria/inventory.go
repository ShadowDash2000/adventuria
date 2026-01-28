package adventuria

import "adventuria/pkg/cache"

type Inventory interface {
	cache.Closable

	MaxSlots() int
	SetMaxSlots(int)
	AvailableSlots() int
	HasEmptySlots() bool
	AddItem(ItemRecord) (string, error)
	AddItemById(string) (string, error)
	MustAddItemById(string) (string, error)
	CanUseItem(string) bool
	UseItem(string) (OnUseSuccess, OnUseFail, error)
	DropItem(string) error
	DropRandomItem() error
	DropInventory() error
	GetItemById(invItemId string) (Item, bool)
}

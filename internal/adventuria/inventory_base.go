package adventuria

import (
	"adventuria/pkg/helper"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type InventoryBase struct {
	locator  ServiceLocator
	user     User
	items    map[string]Item
	maxSlots int
}

func NewInventory(locator ServiceLocator, user User, maxSlots int) (Inventory, error) {
	i := &InventoryBase{
		locator:  locator,
		user:     user,
		maxSlots: maxSlots,
	}

	err := i.fetchInventory()
	if err != nil {
		return nil, err
	}

	i.bindHooks()

	return i, nil
}

func (i *InventoryBase) bindHooks() {
	i.locator.PocketBase().OnRecordAfterCreateSuccess(TableInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("user") == i.user.ID() {
			i.items[e.Record.Id], _ = NewItemFromInventoryRecord(i.locator, i.user, e.Record)
		}
		return e.Next()
	})
	i.locator.PocketBase().OnRecordAfterDeleteSuccess(TableInventory).BindFunc(func(e *core.RecordEvent) error {
		if item, ok := i.items[e.Record.Id]; ok {
			item.Sleep()
			delete(i.items, e.Record.Id)
		}
		return e.Next()
	})
}

func (i *InventoryBase) fetchInventory() error {
	invItems, err := i.locator.PocketBase().FindRecordsByFilter(
		TableInventory,
		"user.id = {:userId}",
		"-created",
		0,
		0,
		dbx.Params{"userId": i.user.ID()},
	)
	if err != nil {
		return err
	}

	i.items = make(map[string]Item)
	for _, invItem := range invItems {
		i.items[invItem.Id], err = NewItemFromInventoryRecord(i.locator, i.user, invItem)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *InventoryBase) MaxSlots() int {
	return i.maxSlots
}

func (i *InventoryBase) SetMaxSlots(maxSlots int) {
	i.maxSlots = maxSlots
}

func (i *InventoryBase) AvailableSlots() int {
	usedSlots := 0
	for _, item := range i.items {
		if item.IsUsingSlot() {
			usedSlots++
		}
	}
	return i.maxSlots - usedSlots
}

func (i *InventoryBase) HasEmptySlots() bool {
	return i.AvailableSlots() > 0
}

func (i *InventoryBase) AddItem(item ItemRecord) error {
	inventoryCollection, err := i.locator.Collections().Get(TableInventory)
	if err != nil {
		return err
	}

	record := core.NewRecord(inventoryCollection)
	record.Set("user", i.user.ID())
	record.Set("item", item.ID())
	record.Set("isActive", item.IsActiveByDefault())
	err = i.locator.PocketBase().Save(record)
	if err != nil {
		return err
	}

	return nil
}

func (i *InventoryBase) AddItemById(itemId string) error {
	item, ok := i.locator.Items().GetById(itemId)
	if !ok {
		return errors.New("item not found")
	}

	if item.IsUsingSlot() && !i.HasEmptySlots() {
		return errors.New("no available slots")
	}

	err := i.AddItem(item)
	if err != nil {
		return err
	}

	return nil
}

// MustAddItemById
// Note: before an item added, checks if there are some empty slots.
// If not, trys to drop a random item from inventory.
func (i *InventoryBase) MustAddItemById(itemId string) error {
	item, ok := i.locator.Items().GetById(itemId)
	if !ok {
		return errors.New("item not found")
	}

	if item.IsUsingSlot() && !i.HasEmptySlots() {
		err := i.DropRandomItem()
		if err != nil {
			return err
		}
	}

	err := i.AddItem(item)
	if err != nil {
		return err
	}

	return nil
}

func (i *InventoryBase) UseItem(itemId string) error {
	item, ok := i.items[itemId]
	if !ok {
		return errors.New("inventory item not found")
	}

	return item.Use()
}

func (i *InventoryBase) DropItem(invItemId string) error {
	item, ok := i.items[invItemId]
	if !ok {
		return errors.New("inventory item not found")
	}

	return item.Drop()
}

// DropRandomItem
// Note: removes random item that uses a slot and can be dropped
func (i *InventoryBase) DropRandomItem() error {
	var items []Item
	for _, item := range i.items {
		if item.IsUsingSlot() && item.CanDrop() {
			items = append(items, item)
		}
	}

	if len(items) == 0 {
		return nil
	}

	err := helper.RandomItemFromSlice(items).Drop()
	if err != nil {
		return err
	}

	return nil
}

func (i *InventoryBase) DropInventory() error {
	for invItemId := range i.items {
		err := i.DropItem(invItemId)
		if err != nil {
			return err
		}
	}
	return nil
}

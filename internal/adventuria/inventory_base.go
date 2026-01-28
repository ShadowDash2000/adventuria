package adventuria

import (
	"adventuria/pkg/helper"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type InventoryBase struct {
	user     User
	items    map[string]Item
	maxSlots int
	hookIds  []string
}

func NewInventory(user User, maxSlots int) (Inventory, error) {
	i := &InventoryBase{
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
	i.hookIds = make([]string, 3)

	i.hookIds[0] = PocketBase.OnRecordAfterCreateSuccess(CollectionInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("user") == i.user.ID() {
			item, err := NewItemFromInventoryRecord(i.user, e.Record)
			if err != nil {
				return err
			}
			i.items[e.Record.Id] = item

			_, err = i.user.OnAfterItemSave().Trigger(&OnAfterItemSave{
				Item: item,
			})
			if err != nil {
				PocketBase.Logger().Error("Failed to trigger OnAfterItemSave event", "err", err)
			}
		}
		return e.Next()
	})
	i.hookIds[1] = PocketBase.OnRecordAfterDeleteSuccess(CollectionInventory).BindFunc(func(e *core.RecordEvent) error {
		if _, ok := i.items[e.Record.Id]; ok {
			delete(i.items, e.Record.Id)
		}
		return e.Next()
	})
	i.hookIds[2] = PocketBase.OnRecordEnrich(CollectionInventory).BindFunc(func(e *core.RecordEnrichEvent) error {
		if _, ok := i.items[e.Record.Id]; ok {
			e.Record.WithCustomData(true)
			e.Record.Set("can_use", i.CanUseItem(e.Record.Id))
		}
		return e.Next()
	})
}

func (i *InventoryBase) Close() {
	PocketBase.OnRecordAfterCreateSuccess(CollectionInventory).Unbind(i.hookIds[0])
	PocketBase.OnRecordAfterDeleteSuccess(CollectionInventory).Unbind(i.hookIds[1])
	PocketBase.OnRecordEnrich(CollectionInventory).Unbind(i.hookIds[2])
	for _, item := range i.items {
		item.Close()
	}
}

func (i *InventoryBase) fetchInventory() error {
	invItems, err := PocketBase.FindRecordsByFilter(
		CollectionInventory,
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
		i.items[invItem.Id], err = NewItemFromInventoryRecord(i.user, invItem)
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

func (i *InventoryBase) AddItem(item ItemRecord) (string, error) {
	onBeforeItemAddEvent := OnBeforeItemAdd{
		ItemRecord:    item,
		ShouldAddItem: true,
	}

	res, err := i.user.OnBeforeItemAdd().Trigger(&onBeforeItemAddEvent)
	if res != nil && !res.Success {
		return "", errors.New(res.Error)
	}
	if err != nil {
		return "", err
	}

	if !onBeforeItemAddEvent.ShouldAddItem {
		return "", nil
	}

	record := core.NewRecord(GameCollections.Get(CollectionInventory))
	record.Set("user", i.user.ID())
	record.Set("item", item.ID())
	record.Set("isActive", item.IsActiveByDefault())
	err = PocketBase.Save(record)
	if err != nil {
		return "", err
	}

	res, err = i.user.OnAfterItemAdd().Trigger(&OnAfterItemAdd{
		ItemRecord: item,
	})
	if res != nil && !res.Success {
		return "", errors.New(res.Error)
	}
	if err != nil {
		return "", err
	}

	return record.Id, nil
}

func (i *InventoryBase) AddItemById(itemId string) (string, error) {
	item, ok := GameItems.GetById(itemId)
	if !ok {
		return "", errors.New("item not found")
	}

	if item.IsUsingSlot() && !i.HasEmptySlots() {
		return "", errors.New("no available slots")
	}

	return i.AddItem(item)
}

// MustAddItemById
// Note: before an item added, checks if there are some empty slots.
// If not, trys to drop a random item from inventory.
func (i *InventoryBase) MustAddItemById(itemId string) (string, error) {
	item, ok := GameItems.GetById(itemId)
	if !ok {
		return "", errors.New("item not found")
	}

	if item.IsUsingSlot() && !i.HasEmptySlots() {
		err := i.DropRandomItem()
		if err != nil {
			return "", err
		}
	}

	return i.AddItem(item)
}

func (i *InventoryBase) CanUseItem(itemId string) bool {
	item, ok := i.items[itemId]
	if !ok {
		return false
	}

	return item.CanUse()
}

func (i *InventoryBase) UseItem(itemId string) (OnUseSuccess, OnUseFail, error) {
	item, ok := i.items[itemId]
	if !ok {
		return nil, nil, errors.New("inventory item not found")
	}

	if !item.CanUse() {
		return nil, nil, errors.New("item cannot be used")
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

func (i *InventoryBase) GetItemById(invItemId string) (Item, bool) {
	item, ok := i.items[invItemId]
	return item, ok
}

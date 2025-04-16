package adventuria

import (
	"adventuria/pkg/helper"
	"errors"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"maps"
	"slices"
	"sort"
)

type InventoryBase struct {
	userId   string
	invItems map[string]InventoryItem
	maxSlots int
}

func NewInventory(userId string, maxSlots int) (Inventory, error) {
	i := &InventoryBase{
		userId:   userId,
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
	GameApp.OnRecordAfterCreateSuccess(TableInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("user") == i.userId {
			i.invItems[e.Record.Id], _ = NewInventoryItemFromRecord(e.Record)
		}
		return e.Next()
	})
	GameApp.OnRecordAfterUpdateSuccess(TableInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("user") == i.userId {
			i.invItems[e.Record.Id].SetProxyRecord(e.Record)
		}
		return e.Next()
	})
	GameApp.OnRecordAfterDeleteSuccess(TableInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("user") == i.userId {
			delete(i.invItems, e.Record.Id)
		}
		return e.Next()
	})
}

func (i *InventoryBase) fetchInventory() error {
	invItems, err := GameApp.FindRecordsByFilter(
		TableInventory,
		"user.id = {:userId}",
		"-created",
		0,
		0,
		dbx.Params{"userId": i.userId},
	)
	if err != nil {
		return err
	}

	i.invItems = make(map[string]InventoryItem)
	for _, invItem := range invItems {
		i.invItems[invItem.Id], err = NewInventoryItemFromRecord(invItem)
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
	for _, item := range i.invItems {
		if item.IsUsingSlot() {
			usedSlots++
		}
	}
	return i.maxSlots - usedSlots
}

func (i *InventoryBase) HasEmptySlots() bool {
	return i.AvailableSlots() > 0
}

func (i *InventoryBase) AddItem(item Item) error {
	inventoryCollection, err := GameCollections.Get(TableInventory)
	if err != nil {
		return err
	}

	record := core.NewRecord(inventoryCollection)
	record.Set("user", i.userId)
	record.Set("item", item.ID())
	record.Set("isActive", item.IsActiveByDefault())
	err = GameApp.Save(record)
	if err != nil {
		return err
	}

	return nil
}

func (i *InventoryBase) AddItemById(itemId string) error {
	item, ok := GameItems.GetById(itemId)
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
// Note: before an item added, checks if there is some empty slots.
// If not, trys to drop a random item from inventory.
func (i *InventoryBase) MustAddItemById(itemId string) error {
	item, ok := GameItems.GetById(itemId)
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

func (i *InventoryBase) Effects(event EffectUse) (*Effects, map[string][]string, error) {
	keys := slices.Collect(maps.Keys(i.invItems))
	slices.Sort(keys)
	sort.Slice(keys, func(k, j int) bool {
		return i.invItems[keys[k]].Order() < i.invItems[keys[j]].Order()
	})

	effects := NewEffects()

	invItemsEffectsIds := make(map[string][]string)
	for _, invItemId := range keys {
		invItem := i.invItems[invItemId]
		itemEffects := invItem.EffectsByEvent(event)
		if itemEffects == nil {
			continue
		}

		var effectsIds []string
		for _, effect := range itemEffects {
			effects.Add(effect)

			// TODO this shouldn't be here
			//i.gc.Log.Add(i.userId, LogTypeItemEffectApplied, effect.Name())
			effectsIds = append(effectsIds, effect.ID())
		}

		invItemsEffectsIds[invItemId] = effectsIds
	}

	return effects, invItemsEffectsIds, nil
}

func (i *InventoryBase) ApplyEffectsByEvent(event EffectUse) (*Effects, error) {
	effects, invItemsEffectsIds, err := i.Effects(event)
	if err != nil {
		return nil, err
	}

	err = i.ApplyEffects(invItemsEffectsIds)

	return effects, nil
}

func (i *InventoryBase) ApplyEffects(invItemsEffectsIds map[string][]string) error {
	for invItemId, effectsIds := range invItemsEffectsIds {
		invItem := i.invItems[invItemId]

		err := invItem.ApplyEffects(effectsIds)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *InventoryBase) ApplyEffectsByTypes(types []string) error {
	var effectsIds map[string][]string
	for _, invItem := range i.invItems {
		for _, effect := range invItem.EffectsByTypes(types) {
			effectsIds[invItem.ID()] = append(effectsIds[invItem.ID()], effect.ID())
		}
	}

	return i.ApplyEffects(effectsIds)
}

func (i *InventoryBase) UseItem(itemId string) error {
	item, ok := i.invItems[itemId]
	if !ok {
		return errors.New("inventory item not found")
	}

	return item.Use()
}

func (i *InventoryBase) DropItem(invItemId string) error {
	invItem, ok := i.invItems[invItemId]
	if !ok {
		return errors.New("inventory item not found")
	}

	return invItem.Drop()
}

// DropRandomItem
// Note: removes random item that uses a slot and can be dropped
func (i *InventoryBase) DropRandomItem() error {
	var items []InventoryItem
	for _, item := range i.invItems {
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
	for invItemId, _ := range i.invItems {
		err := i.DropItem(invItemId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *InventoryBase) Items() map[string]InventoryItem {
	return i.invItems
}

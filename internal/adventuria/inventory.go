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

type Inventory struct {
	gc       *GameComponents
	userId   string
	invItems map[string]*InventoryItem
	maxSlots int
}

func NewInventory(userId string, maxSlots int, gc *GameComponents) (*Inventory, error) {
	i := &Inventory{
		gc:       gc,
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

func (i *Inventory) bindHooks() {
	i.gc.App.OnRecordAfterCreateSuccess(TableInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("user") == i.userId {
			i.invItems[e.Record.Id], _ = NewInventoryItem(e.Record, i.gc)
		}
		return e.Next()
	})
	i.gc.App.OnRecordAfterUpdateSuccess(TableInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("user") == i.userId {
			i.invItems[e.Record.Id].SetProxyRecord(e.Record)
		}
		return e.Next()
	})
	i.gc.App.OnRecordAfterDeleteSuccess(TableInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("user") == i.userId {
			delete(i.invItems, e.Record.Id)
		}
		return e.Next()
	})
}

func (i *Inventory) fetchInventory() error {
	invItems, err := i.gc.App.FindRecordsByFilter(
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

	i.invItems = make(map[string]*InventoryItem)
	for _, invItem := range invItems {
		i.invItems[invItem.Id], err = NewInventoryItem(invItem, i.gc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Inventory) SetMaxSlots(maxSlots int) {
	i.maxSlots = maxSlots
}

func (i *Inventory) GetAvailableSlots() int {
	usedSlots := 0
	for _, item := range i.invItems {
		if item.IsUsingSlot() {
			usedSlots++
		}
	}
	return i.maxSlots - usedSlots
}

func (i *Inventory) HasEmptySlots() bool {
	return i.GetAvailableSlots() > 0
}

func (i *Inventory) AddItem(itemId string) error {
	item, err := i.gc.App.FindRecordById(TableItems, itemId)
	if err != nil {
		return err
	}

	if item.GetBool("isUsingSlot") && !i.HasEmptySlots() {
		return errors.New("no available slots")
	}

	err = i.CreateInventoryRecord(item)
	if err != nil {
		return err
	}

	return nil
}

// MustAddItem
// Note: before an item added, checks if there is some empty slots.
// If not, trys to drop a random item from inventory.
func (i *Inventory) MustAddItem(itemId string) error {
	item, err := i.gc.App.FindRecordById(TableItems, itemId)
	if err != nil {
		return err
	}

	if item.GetBool("isUsingSlot") && !i.HasEmptySlots() {
		err := i.DropRandomItem()
		if err != nil {
			return err
		}
	}

	err = i.CreateInventoryRecord(item)
	if err != nil {
		return err
	}
	return nil
}

func (i *Inventory) CreateInventoryRecord(item *core.Record) error {
	inventoryCollection, err := i.gc.Cols.Get(TableInventory)
	if err != nil {
		return err
	}

	record := core.NewRecord(inventoryCollection)
	record.Set("user", i.userId)
	record.Set("item", item.Id)
	record.Set("isActive", item.GetBool("isActiveByDefault"))
	err = i.gc.App.Save(record)
	if err != nil {
		return err
	}

	return nil
}

func (i *Inventory) GetEffects(event string) (*Effects, map[string][]string, error) {
	keys := slices.Collect(maps.Keys(i.invItems))
	slices.Sort(keys)
	sort.Slice(keys, func(k, j int) bool {
		return i.invItems[keys[k]].item.Order() < i.invItems[keys[j]].item.Order()
	})

	effects := NewEffects()

	invItemsEffectsIds := make(map[string][]string)
	for _, invItemId := range keys {
		invItem := i.invItems[invItemId]
		itemEffects := invItem.GetEffectsByEvent(event)
		if itemEffects == nil {
			continue
		}

		var effectsIds []string
		for _, effect := range itemEffects {
			effects.AddValue(effect.Type(), effect.Value())

			i.gc.Log.Add(i.userId, LogTypeItemEffectApplied, effect.Name())
			effectsIds = append(effectsIds, effect.GetId())
		}

		invItemsEffectsIds[invItemId] = effectsIds
	}

	return effects, invItemsEffectsIds, nil
}

func (i *Inventory) applyEffects(event string) (*Effects, error) {
	effects, invItemsEffectsIds, err := i.GetEffects(event)
	if err != nil {
		return nil, err
	}

	for invItemId, effectsIds := range invItemsEffectsIds {
		invItem := i.invItems[invItemId]

		appliedEffects := invItem.AppliedEffects()
		invItem.SetAppliedEffects(append(appliedEffects, effectsIds...))
		appliedEffectsCount := len(invItem.AppliedEffects())

		if appliedEffectsCount < invItem.EffectsCount() {
			err = i.gc.App.Save(invItem)
			if err != nil {
				return nil, err
			}
		} else {
			err = i.gc.App.Delete(invItem)
			if err != nil {
				return nil, err
			}
		}
	}

	return effects, nil
}

func (i *Inventory) UseItem(itemId string) error {
	item, ok := i.invItems[itemId]
	if !ok {
		return errors.New("item not found")
	}

	i.gc.Log.Add(i.userId, LogTypeItemUse, item.GetName())

	return item.Use()
}

func (i *Inventory) DropItem(invItemId string) error {
	invItem, ok := i.invItems[invItemId]
	if !ok {
		return errors.New("inventory item not found")
	}

	if !invItem.CanDrop() {
		return errors.New("inventory item isn't droppable")
	}

	err := i.gc.App.Delete(invItem)
	if err != nil {
		return err
	}

	i.gc.Log.Add(i.userId, LogTypeItemDrop, invItem.GetName())

	return nil
}

// DropRandomItem
// Note: removes random item that uses a slot and can be dropped
func (i *Inventory) DropRandomItem() error {
	var items []*InventoryItem
	for _, item := range i.invItems {
		if item.IsUsingSlot() && item.CanDrop() {
			items = append(items, item)
		}
	}

	if len(items) == 0 {
		return nil
	}

	err := i.DropItem(helper.RandomItemFromSlice(items).Id)
	if err != nil {
		return err
	}

	return nil
}

func (i *Inventory) DropInventory() error {
	for invItemId, _ := range i.invItems {
		err := i.DropItem(invItemId)
		if err != nil {
			return err
		}
	}
	return nil
}

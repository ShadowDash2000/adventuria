package adventuria

import (
	"adventuria/pkg/collections"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"maps"
	"slices"
	"sort"
)

type Inventory struct {
	app      core.App
	cols     *collections.Collections
	log      *Log
	userId   string
	invItems map[string]*InventoryItem
	maxSlots int
}

func NewInventory(userId string, maxSlots int, log *Log, cols *collections.Collections, app core.App) (*Inventory, error) {
	i := &Inventory{
		app:      app,
		cols:     cols,
		log:      log,
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
	i.app.OnRecordAfterCreateSuccess(TableInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("user") == i.userId {
			i.invItems[e.Record.Id], _ = NewInventoryItem(e.Record, i.log, i.app)
		}
		return e.Next()
	})
	i.app.OnRecordAfterUpdateSuccess(TableInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("user") == i.userId {
			i.invItems[e.Record.Id].SetProxyRecord(e.Record)
		}
		return e.Next()
	})
	i.app.OnRecordAfterDeleteSuccess(TableInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("user") == i.userId {
			delete(i.invItems, e.Record.Id)
		}
		return e.Next()
	})
}

func (i *Inventory) fetchInventory() error {
	inventory, err := i.app.FindRecordsByFilter(
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
	for _, item := range inventory {
		i.invItems[item.Id], err = NewInventoryItem(item, i.log, i.app)
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

func (i *Inventory) AddItem(itemId string) error {
	if i.GetAvailableSlots() <= 0 {
		return errors.New("no available slots")
	}

	inventoryCollection, err := i.cols.Get(TableInventory)
	if err != nil {
		return err
	}

	item, err := i.app.FindRecordById(TableItems, itemId)
	if err != nil {
		return err
	}

	record := core.NewRecord(inventoryCollection)
	record.Set("user", i.userId)
	record.Set("item", itemId)
	record.Set("isActive", item.GetBool("isActiveByDefault"))
	err = i.app.Save(record)
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

	var effects *Effects
	var effectsMap map[string]interface{}

	err := mapstructure.Decode(&Effects{}, &effectsMap)
	if err != nil {
		return nil, nil, err
	}

	invItemsEffectsIds := make(map[string][]string)
	for _, invItemId := range keys {
		invItem := i.invItems[invItemId]
		itemEffects := invItem.GetEffectsByEvent(event)
		if itemEffects == nil {
			continue
		}

		var effectsIds []string
		for _, effect := range itemEffects {
			switch effect.Kind() {
			case Int:
				effectsMap[effect.Type()] = effectsMap[effect.Type()].(int) + effect.GetInt()
			case Bool:
				effectsMap[effect.Type()] = true
			case Slice:
				effectsMap[effect.Type()] = effect.GetSlice()
			}

			effectsIds = append(effectsIds, effect.Id())
		}

		invItemsEffectsIds[invItemId] = effectsIds
	}

	err = mapstructure.Decode(effectsMap, &effects)
	if err != nil {
		return nil, nil, err
	}

	return effects, invItemsEffectsIds, nil
}

func (i *Inventory) ApplyEffects(event string) (*Effects, error) {
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
			err = i.app.Save(invItem)
			if err != nil {
				return nil, err
			}
		} else {
			err = i.app.Delete(invItem)
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

	i.log.Add(i.userId, LogTypeItemUse, item.GetName())

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

	err := i.app.Delete(invItem)
	if err != nil {
		return err
	}

	i.log.Add(i.userId, LogTypeItemDrop, invItem.GetName())

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

package adventuria

import (
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
	userId   string
	items    map[string]*InventoryItem
	maxSlots int
}

func NewInventory(userId string, maxSlots int, app core.App) (*Inventory, error) {
	i := &Inventory{
		app:      app,
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
			i.items[e.Record.Id], _ = NewInventoryItem(e.Record, i.app)
		}
		return e.Next()
	})
	i.app.OnRecordAfterUpdateSuccess(TableInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("user") == i.userId {
			i.items[e.Record.Id].invItem = e.Record
		}
		return e.Next()
	})
	i.app.OnRecordAfterDeleteSuccess(TableInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("user") == i.userId {
			delete(i.items, e.Record.Id)
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

	i.items = make(map[string]*InventoryItem)
	for _, item := range inventory {
		i.items[item.Id], err = NewInventoryItem(item, i.app)
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
	for _, item := range i.items {
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

	inventoryCollection, err := i.app.FindCollectionByNameOrId(TableInventory)
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

func (i *Inventory) GetEffects(event string) (*Effects, []string, error) {
	keys := slices.Collect(maps.Keys(i.items))
	slices.Sort(keys)
	sort.Slice(keys, func(k, j int) bool {
		return i.items[keys[k]].item.GetOrder() < i.items[keys[j]].item.GetOrder()
	})

	var effects *Effects
	var effectsMap map[string]interface{}

	err := mapstructure.Decode(&Effects{}, &effectsMap)
	if err != nil {
		return nil, nil, err
	}

	var itemsIds []string
	for _, itemId := range keys {
		item := i.items[itemId]
		itemEffects := item.GetEffects(event)
		if itemEffects == nil {
			continue
		}

		for _, effect := range itemEffects {
			switch effect.Kind() {
			case Int:
				effectsMap[effect.Type()] = effectsMap[effect.Type()].(int) + effect.GetInt()
			case Bool:
				effectsMap[effect.Type()] = true
			case Slice:
				effectsMap[effect.Type()] = effect.GetSlice()
			}
		}

		itemsIds = append(itemsIds, item.invItem.Id)
	}

	err = mapstructure.Decode(effectsMap, &effects)
	if err != nil {
		return nil, nil, err
	}

	return effects, itemsIds, nil
}

func (i *Inventory) ApplyEffects(event string) (*Effects, error) {
	effects, itemsIds, err := i.GetEffects(event)
	if err != nil {
		return nil, err
	}

	if len(itemsIds) > 0 {
		err := i.app.RunInTransaction(func(txApp core.App) error {
			for _, itemId := range itemsIds {
				err := txApp.Delete(i.items[itemId].invItem)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return effects, nil
}

package adventuria

import (
	"errors"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"maps"
	"reflect"
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

func (i *Inventory) GetEffects(event string) (any, []string) {
	var effects any

	switch event {
	case ItemUseTypeOnDrop:
		effects = &OnDropEffects{}
	case ItemUseTypeOnRoll:
		effects = &OnRollEffects{}
	default:
		return nil, nil
	}

	keys := slices.Collect(maps.Keys(i.items))
	slices.Sort(keys)
	sort.Slice(keys, func(k, j int) bool {
		return i.items[keys[k]].item.GetOrder() < i.items[keys[j]].item.GetOrder()
	})

	var itemsIds []string
	effectsValue := reflect.ValueOf(effects).Elem()
	fields := reflect.VisibleFields(effectsValue.Type())
	for _, itemId := range keys {
		item := i.items[itemId]
		itemEffects := item.GetEffects(event)
		if itemEffects == nil {
			continue
		}

		itemValue := reflect.ValueOf(itemEffects).Elem()
		for _, field := range fields {
			fieldValue := itemValue.FieldByIndex(field.Index)
			if fieldValue.IsZero() {
				continue
			}

			switch field.Type.Kind() {
			case reflect.Int:
				v1 := effectsValue.FieldByIndex(field.Index).Int()
				v2 := fieldValue.Int()
				effectsValue.FieldByIndex(field.Index).SetInt(v1 + v2)
			case reflect.Bool:
				effectsValue.FieldByIndex(field.Index).Set(fieldValue)
			case reflect.Slice:
				if effectsValue.FieldByIndex(field.Index).IsZero() {
					effectsValue.FieldByIndex(field.Index).Set(fieldValue)
				}
			default:
			}
		}

		itemsIds = append(itemsIds, item.invItem.Id)
	}

	return effects, itemsIds
}

func (i *Inventory) ApplyEffects(event string) (any, error) {
	effects, itemsIds := i.GetEffects(event)

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

func (i *Inventory) GetOnRollEffects() OnRollEffects {
	effects, _ := i.GetEffects(ItemUseTypeOnRoll)
	if effects == nil {
		return OnRollEffects{}
	}

	return *(effects.(*OnRollEffects))
}

func (i *Inventory) ApplyOnDropEffects() (OnDropEffects, error) {
	effects, err := i.ApplyEffects(ItemUseTypeOnDrop)
	if err != nil {
		return OnDropEffects{}, err
	}
	if effects == nil {
		return OnDropEffects{}, nil
	}

	return *(effects.(*OnDropEffects)), nil
}

func (i *Inventory) ApplyOnRollEffects() (OnRollEffects, error) {
	effects, err := i.ApplyEffects(ItemUseTypeOnRoll)
	if err != nil {
		return OnRollEffects{}, err
	}
	if effects == nil {
		return OnRollEffects{}, nil
	}

	return *(effects.(*OnRollEffects)), nil
}

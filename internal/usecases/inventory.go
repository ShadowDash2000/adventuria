package usecases

import (
	"adventuria/internal/adventuria"
	"errors"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

var inventoryCollection *core.Collection
var itemsList map[string]*core.Record

type Inventory struct {
	app       core.App
	userId    string
	inventory map[string]*core.Record
	items     map[string]Usable
}

func NewInventory(userId string, app core.App) (*Inventory, error) {
	var err error
	i := &Inventory{
		app:    app,
		userId: userId,
	}

	if inventoryCollection == nil {
		inventoryCollection, err = app.FindCollectionByNameOrId(adventuria.TableInventory)
		if err != nil {
			return nil, err
		}
	}

	if itemsList == nil {
		items, err := app.FindRecordsByFilter(
			adventuria.TableItems,
			"",
			"",
			0,
			0,
		)
		if err != nil {
			return nil, err
		}

		itemsList = make(map[string]*core.Record, len(items))
		for _, item := range items {
			itemsList[item.Id] = item
		}
	}

	inventory, err := app.FindRecordsByFilter(
		adventuria.TableInventory,
		"user.id = {:userId}",
		"-created",
		0,
		0,
		dbx.Params{"userId": userId},
	)
	if err != nil {
		return nil, err
	}

	i.inventory = make(map[string]*core.Record)
	i.items = make(map[string]Usable)
	for _, record := range inventory {
		recordFields := record.FieldsData()
		itemFields := itemsList[recordFields["item"].(string)].FieldsData()

		i.inventory[record.Id] = record
		i.items[record.Id] = NewItem(itemFields["type"].(string), userId)
	}

	return i, nil
}

func (i *Inventory) AddItemById(itemId string) error {
	_, ok := itemsList[itemId]
	if !ok {
		return errors.New("item not found")
	}

	record := core.NewRecord(inventoryCollection)
	record.Set("user", i.userId)
	record.Set("item", itemId)
	err := i.app.Save(record)
	if err != nil {
		return err
	}

	itemFields := itemsList[itemId].FieldsData()
	i.inventory[record.Id] = record
	i.items[record.Id] = NewItem(itemFields["type"].(string), i.userId)

	return nil
}

func (i *Inventory) UseItemById(invItemId string) error {
	_, ok := i.inventory[invItemId]
	if !ok {
		return errors.New("item not found in inventory")
	}

	item, ok := i.items[invItemId]
	if !ok {
		return errors.New("item not found in inventory")
	}

	err := item.Use()
	if err != nil {
		return err
	}

	err = i.app.Delete(i.inventory[invItemId])
	if err != nil {
		return err
	}

	delete(i.inventory, invItemId)
	delete(i.items, invItemId)

	return nil
}

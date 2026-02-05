package adventuria

import (
	"adventuria/pkg/helper"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type InventoryBase struct {
	user     User
	items    map[string]Item
	maxSlots int
	hookIds  []string
}

func NewInventory(ctx AppContext, user User, maxSlots int) (Inventory, error) {
	i := &InventoryBase{
		user:     user,
		maxSlots: maxSlots,
	}

	err := i.fetchInventory(ctx)
	if err != nil {
		return nil, err
	}

	i.bindHooks(ctx)

	return i, nil
}

func (i *InventoryBase) bindHooks(ctx AppContext) {
	i.hookIds = make([]string, 4)

	i.hookIds[0] = ctx.App.OnRecordAfterCreateSuccess(CollectionInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("user") == i.user.ID() {
			ctx := AppContext{App: e.App}
			item, err := NewItemFromInventoryRecord(ctx, i.user, e.Record)
			if err != nil {
				return err
			}
			i.items[e.Record.Id] = item

			_, err = i.user.OnAfterItemSave().Trigger(&OnAfterItemSave{
				AppContext: ctx,
				Item:       item,
			})
			if err != nil {
				e.App.Logger().Error("Failed to trigger OnAfterItemSave event", "err", err)
			}
		}
		return e.Next()
	})
	i.hookIds[1] = ctx.App.OnRecordAfterDeleteSuccess(CollectionInventory).BindFunc(func(e *core.RecordEvent) error {
		if _, ok := i.items[e.Record.Id]; ok {
			delete(i.items, e.Record.Id)
		}
		return e.Next()
	})
	i.hookIds[2] = ctx.App.OnRecordEnrich(CollectionInventory).BindFunc(func(e *core.RecordEnrichEvent) error {
		if _, ok := i.items[e.Record.Id]; ok {
			e.Record.WithCustomData(true)
			e.Record.Set("can_use", i.CanUseItem(AppContext{App: e.App}, e.Record.Id))
		}
		return e.Next()
	})
	i.hookIds[3] = ctx.App.OnRecordCreate(CollectionInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetBool("isActive") {
			e.Record.Set("activated", types.NowDateTime())
		}
		return e.Next()
	})
}

func (i *InventoryBase) Close(ctx AppContext) {
	ctx.App.OnRecordAfterCreateSuccess(CollectionInventory).Unbind(i.hookIds[0])
	ctx.App.OnRecordAfterDeleteSuccess(CollectionInventory).Unbind(i.hookIds[1])
	ctx.App.OnRecordEnrich(CollectionInventory).Unbind(i.hookIds[2])
	ctx.App.OnRecordCreate(CollectionInventory).Unbind(i.hookIds[3])
	for _, item := range i.items {
		item.Close(ctx)
	}
}

func (i *InventoryBase) fetchInventory(ctx AppContext) error {
	invItems, err := ctx.App.FindRecordsByFilter(
		CollectionInventory,
		"user.id = {:userId}",
		"activated,created",
		0,
		0,
		dbx.Params{"userId": i.user.ID()},
	)
	if err != nil {
		return err
	}

	i.items = make(map[string]Item)
	for _, invItem := range invItems {
		i.items[invItem.Id], err = NewItemFromInventoryRecord(ctx, i.user, invItem)
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

func (i *InventoryBase) AddItem(ctx AppContext, item ItemRecord) (string, error) {
	onBeforeItemAddEvent := OnBeforeItemAdd{
		AppContext:    ctx,
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
	err = ctx.App.Save(record)
	if err != nil {
		return "", err
	}

	res, err = i.user.OnAfterItemAdd().Trigger(&OnAfterItemAdd{
		AppContext: ctx,
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

func (i *InventoryBase) AddItemById(ctx AppContext, itemId string) (string, error) {
	item, ok := GameItems.GetById(itemId)
	if !ok {
		return "", errors.New("item not found")
	}

	if item.IsUsingSlot() && !i.HasEmptySlots() {
		return "", errors.New("no available slots")
	}

	return i.AddItem(ctx, item)
}

// MustAddItemById
// Note: before an item added, checks if there are some empty slots.
// If not, trys to drop a random item from inventory.
func (i *InventoryBase) MustAddItemById(ctx AppContext, itemId string) (string, error) {
	item, ok := GameItems.GetById(itemId)
	if !ok {
		return "", errors.New("item not found")
	}

	if item.IsUsingSlot() && !i.HasEmptySlots() {
		err := i.DropRandomItem(ctx)
		if err != nil {
			return "", err
		}
	}

	return i.AddItem(ctx, item)
}

func (i *InventoryBase) CanUseItem(ctx AppContext, itemId string) bool {
	item, ok := i.items[itemId]
	if !ok {
		return false
	}

	return item.CanUse(ctx)
}

func (i *InventoryBase) UseItem(ctx AppContext, itemId string) (OnUseSuccess, OnUseFail, error) {
	item, ok := i.items[itemId]
	if !ok {
		return nil, nil, errors.New("inventory item not found")
	}

	if !item.CanUse(ctx) {
		return nil, nil, errors.New("item cannot be used")
	}

	return item.Use(ctx)
}

func (i *InventoryBase) DropItem(ctx AppContext, invItemId string) error {
	item, ok := i.items[invItemId]
	if !ok {
		return errors.New("inventory item not found")
	}

	return item.Drop(ctx)
}

// DropRandomItem
// Note: removes random item that uses a slot and can be dropped
func (i *InventoryBase) DropRandomItem(ctx AppContext) error {
	var items []Item
	for _, item := range i.items {
		if item.IsUsingSlot() && item.CanDrop() {
			items = append(items, item)
		}
	}

	if len(items) == 0 {
		return nil
	}

	err := helper.RandomItemFromSlice(items).Drop(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (i *InventoryBase) DropInventory(ctx AppContext) error {
	for invItemId := range i.items {
		err := i.DropItem(ctx, invItemId)
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

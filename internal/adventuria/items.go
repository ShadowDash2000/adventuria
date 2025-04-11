package adventuria

import (
	"adventuria/pkg/cache"
	"github.com/pocketbase/pocketbase/core"
)

type Items struct {
	gc    *GameComponents
	items *cache.MemoryCache[string, Item]
}

func NewItems(gc *GameComponents) *Items {
	items := &Items{
		gc:    gc,
		items: cache.NewMemoryCache[string, Item](0, true),
	}

	items.fetch()
	items.bindHooks()

	return items
}

func (i *Items) bindHooks() {
	i.gc.App.OnRecordAfterCreateSuccess(TableItems).BindFunc(func(e *core.RecordEvent) error {
		i.add(e.Record)
		return e.Next()
	})
	i.gc.App.OnRecordAfterUpdateSuccess(TableItems).BindFunc(func(e *core.RecordEvent) error {
		i.add(e.Record)
		return e.Next()
	})
	i.gc.App.OnRecordAfterDeleteSuccess(TableItems).BindFunc(func(e *core.RecordEvent) error {
		i.delete(e.Record.Id)
		return e.Next()
	})
}

func (i *Items) fetch() error {
	i.items.Clear()

	items, err := i.gc.App.FindAllRecords(TableItems)
	if err != nil {
		return err
	}

	for _, item := range items {
		err = i.add(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Items) add(record *core.Record) error {
	item, err := NewItemFromRecord(record, i.gc)
	if err != nil {
		return err
	}

	i.items.Set(item.ID(), item)

	return nil
}

func (i *Items) delete(id string) {
	i.items.Delete(id)
}

func (i *Items) GetById(id string) (Item, bool) {
	return i.items.Get(id)
}

func (i *Items) GetAll() map[string]Item {
	return i.items.GetAll()
}

func (i *Items) GetAllRollable() []Item {
	var res []Item
	for _, item := range i.items.GetAll() {
		if !item.IsRollable() {
			continue
		}
		res = append(res, item)
	}
	return res
}

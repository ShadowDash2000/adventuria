package adventuria

import (
	"adventuria/pkg/cache"
	"github.com/pocketbase/pocketbase/core"
)

type Items struct {
	items *cache.MemoryCache[string, Item]
}

func NewItems() *Items {
	items := &Items{
		items: cache.NewMemoryCache[string, Item](0, true),
	}

	items.fetch()
	items.bindHooks()

	return items
}

func (i *Items) bindHooks() {
	GameApp.OnRecordAfterCreateSuccess(TableItems).BindFunc(func(e *core.RecordEvent) error {
		i.add(e.Record)
		return e.Next()
	})
	GameApp.OnRecordAfterUpdateSuccess(TableItems).BindFunc(func(e *core.RecordEvent) error {
		i.add(e.Record)
		return e.Next()
	})
	GameApp.OnRecordAfterDeleteSuccess(TableItems).BindFunc(func(e *core.RecordEvent) error {
		i.delete(e.Record.Id)
		return e.Next()
	})
}

func (i *Items) fetch() error {
	i.items.Clear()

	items, err := GameApp.FindAllRecords(TableItems)
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
	item, err := NewItemFromRecord(record)
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

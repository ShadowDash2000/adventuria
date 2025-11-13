package adventuria

import (
	"adventuria/pkg/cache"
	"iter"

	"github.com/pocketbase/pocketbase/core"
)

type Items struct {
	items *cache.MemoryCache[string, ItemRecord]
}

func NewItems() *Items {
	items := &Items{
		items: cache.NewMemoryCache[string, ItemRecord](0, true),
	}

	items.fetch()
	items.bindHooks()

	return items
}

func (i *Items) bindHooks() {
	PocketBase.OnRecordAfterCreateSuccess(CollectionItems).BindFunc(func(e *core.RecordEvent) error {
		i.add(e.Record)
		return e.Next()
	})
	PocketBase.OnRecordAfterUpdateSuccess(CollectionItems).BindFunc(func(e *core.RecordEvent) error {
		i.add(e.Record)
		return e.Next()
	})
	PocketBase.OnRecordAfterDeleteSuccess(CollectionItems).BindFunc(func(e *core.RecordEvent) error {
		i.delete(e.Record.Id)
		return e.Next()
	})
}

func (i *Items) fetch() error {
	i.items.Clear()

	items, err := PocketBase.FindAllRecords(CollectionItems)
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
	item := NewItemFromRecord(record)

	i.items.Set(item.ID(), item)

	return nil
}

func (i *Items) delete(id string) {
	i.items.Delete(id)
}

func (i *Items) GetById(id string) (ItemRecord, bool) {
	return i.items.Get(id)
}

func (i *Items) GetAll() iter.Seq2[string, ItemRecord] {
	return i.items.GetAll()
}

func (i *Items) GetAllRollable() []ItemRecord {
	var res []ItemRecord
	for _, item := range i.items.GetAll() {
		if !item.IsRollable() {
			continue
		}
		res = append(res, item)
	}
	return res
}

package adventuria

import (
	"adventuria/pkg/cache"
	"iter"

	"github.com/pocketbase/pocketbase/core"
)

type Items struct {
	items cache.Cache[string, ItemRecord]
}

func NewItems(ctx AppContext) (*Items, error) {
	items := &Items{
		items: cache.NewMemoryCache[string, ItemRecord](0, true),
	}

	if err := items.fetch(ctx); err != nil {
		return nil, err
	}
	items.bindHooks(ctx)

	return items, nil
}

func (i *Items) bindHooks(ctx AppContext) {
	ctx.App.OnRecordAfterCreateSuccess(CollectionItems).BindFunc(func(e *core.RecordEvent) error {
		i.add(e.Record)
		return e.Next()
	})
	ctx.App.OnRecordAfterUpdateSuccess(CollectionItems).BindFunc(func(e *core.RecordEvent) error {
		i.add(e.Record)
		return e.Next()
	})
	ctx.App.OnRecordAfterDeleteSuccess(CollectionItems).BindFunc(func(e *core.RecordEvent) error {
		i.delete(e.Record.Id)
		return e.Next()
	})
}

func (i *Items) fetch(ctx AppContext) error {
	i.items.Clear()

	items, err := ctx.App.FindAllRecords(CollectionItems)
	if err != nil {
		return err
	}

	for _, item := range items {
		if err = i.add(item); err != nil {
			ctx.App.Logger().Error("Items: unknown item type", "item", item)
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

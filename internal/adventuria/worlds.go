package adventuria

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/cache"
	"database/sql"
	"errors"
	"iter"
	"sync"

	"github.com/pocketbase/pocketbase/core"
)

type Worlds struct {
	worlds cache.Cache[string, World]
	mx     sync.Mutex
}

func NewWorlds(ctx AppContext) (*Worlds, error) {
	w := &Worlds{
		worlds: cache.NewMemoryCache[string, World](0, true),
	}

	if err := w.fetch(ctx); err != nil {
		return nil, err
	}
	w.bindHooks(ctx)

	return w, nil
}

func (w *Worlds) bindHooks(ctx AppContext) {
	ctx.App.OnRecordAfterCreateSuccess(schema.CollectionsWorlds).BindFunc(func(e *core.RecordEvent) error {
		w.add(e.Record)
		return e.Next()
	})
	ctx.App.OnRecordAfterUpdateSuccess(schema.CollectionsWorlds).BindFunc(func(e *core.RecordEvent) error {
		w.add(e.Record)
		return e.Next()
	})
	ctx.App.OnRecordAfterDeleteSuccess(schema.CollectionsWorlds).BindFunc(func(e *core.RecordEvent) error {
		w.delete(e.Record)
		return e.Next()
	})
}

func (w *Worlds) fetch(ctx AppContext) error {
	w.worlds.Clear()

	var worlds []*core.Record
	err := ctx.App.RecordQuery(schema.CollectionsWorlds).
		OrderBy(schema.WorldsSchema.Sort).
		All(&worlds)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	for _, record := range worlds {
		w.add(record)
	}

	return nil
}

func (w *Worlds) add(record *core.Record) {
	w.mx.Lock()
	defer w.mx.Unlock()
	w.worlds.Set(record.Id, NewWorld(record))
}

func (w *Worlds) delete(record *core.Record) {
	w.mx.Lock()
	defer w.mx.Unlock()
	w.worlds.Delete(record.Id)
}

func (w *Worlds) GetByID(id string) (World, bool) {
	return w.worlds.Get(id)
}

func (w *Worlds) GetAll() iter.Seq2[string, World] {
	return w.worlds.GetAll()
}

func (w *Worlds) GetDefault() (World, bool) {
	for _, world := range w.GetAll() {
		if world.IsDefaultWorld() {
			return world, true
		}
	}
	return nil, false
}

package collections

import (
	"adventuria/pkg/cache"

	"github.com/pocketbase/pocketbase/core"
)

type Collections struct {
	app  core.App
	cols *cache.MemoryCache[string, *core.Collection]
}

func NewCollections(app core.App) *Collections {
	return &Collections{
		app:  app,
		cols: cache.NewMemoryCache[string, *core.Collection](0, true),
	}
}

func (c *Collections) Get(name string) *core.Collection {
	if col, ok := c.cols.Get(name); ok {
		return col
	}

	col, err := c.fetchCollection(name)
	if err != nil {
		panic(err)
	}

	c.cols.Set(name, col)
	return col
}

func (c *Collections) fetchCollection(name string) (*core.Collection, error) {
	return c.app.FindCollectionByNameOrId(name)
}

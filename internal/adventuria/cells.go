package adventuria

import (
	"adventuria/pkg/cache"
	"github.com/pocketbase/pocketbase/core"
)

type Cells struct {
	app         core.App
	cellsBySort *cache.MemoryCache[int, *core.Record]
	cellsByCode *cache.MemoryCache[string, *core.Record]
}

func NewCells(app core.App) *Cells {
	cells := &Cells{
		app:         app,
		cellsBySort: cache.NewMemoryCache[int, *core.Record](0, true),
		cellsByCode: cache.NewMemoryCache[string, *core.Record](0, true),
	}

	cells.fetch()
	cells.bindHooks()

	return cells
}

func (c *Cells) bindHooks() {
	c.app.OnRecordAfterCreateSuccess(TableCells).BindFunc(func(e *core.RecordEvent) error {
		c.cellsBySort.Set(e.Record.GetInt("sort"), e.Record)
		if cellCode := e.Record.GetString("code"); cellCode != "" {
			c.cellsByCode.Set(cellCode, e.Record)
		}
		return e.Next()
	})
	c.app.OnRecordAfterUpdateSuccess(TableCells).BindFunc(func(e *core.RecordEvent) error {
		c.cellsBySort.Set(e.Record.GetInt("sort"), e.Record)
		if cellCode := e.Record.GetString("code"); cellCode != "" {
			c.cellsByCode.Set(cellCode, e.Record)
		}
		return e.Next()
	})
	c.app.OnRecordAfterDeleteSuccess(TableCells).BindFunc(func(e *core.RecordEvent) error {
		c.cellsBySort.Delete(e.Record.GetInt("sort"))
		if cellCode := e.Record.GetString("code"); cellCode != "" {
			c.cellsByCode.Delete(cellCode)
		}
		return e.Next()
	})
}

func (c *Cells) fetch() error {
	c.cellsBySort.Clear()
	c.cellsByCode.Clear()

	cells, err := c.app.FindRecordsByFilter(
		TableCells,
		"",
		"sort",
		0,
		0,
	)
	if err != nil {
		return err
	}

	for _, cell := range cells {
		c.cellsBySort.Set(cell.GetInt("sort"), cell)

		code := cell.GetString("code")
		if code != "" {
			c.cellsByCode.Set(code, cell)
		}
	}

	return nil
}

func (c *Cells) GetBySort(sort int) (*core.Record, bool) {
	return c.cellsBySort.Get(sort)
}

func (c *Cells) GetByCode(code string) (*core.Record, bool) {
	return c.cellsByCode.Get(code)
}

func (c *Cells) GetAll() map[int]*core.Record {
	return c.cellsBySort.GetAll()
}

func (c *Cells) CellsByCode() *cache.MemoryCache[string, *core.Record] {
	return c.cellsByCode
}

func (c *Cells) Count() int {
	return c.cellsBySort.Count()
}

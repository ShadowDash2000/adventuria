package adventuria

import (
	"adventuria/pkg/cache"
	"github.com/pocketbase/pocketbase/core"
	"sort"
	"sync"
)

type Cells struct {
	app           core.App
	cells         *cache.MemoryCache[string, *core.Record]
	cellsByCode   *cache.MemoryCache[string, *core.Record]
	cellsOrder    []string
	cellsIdsOrder map[string]int
	mx            sync.Mutex
}

func NewCells(app core.App) *Cells {
	cells := &Cells{
		app:         app,
		cells:       cache.NewMemoryCache[string, *core.Record](0, true),
		cellsByCode: cache.NewMemoryCache[string, *core.Record](0, true),
	}

	cells.fetch()
	cells.bindHooks()

	return cells
}

func (c *Cells) bindHooks() {
	c.app.OnRecordAfterCreateSuccess(TableCells).BindFunc(func(e *core.RecordEvent) error {
		c.add(e.Record)
		c.sort()
		return e.Next()
	})
	c.app.OnRecordAfterUpdateSuccess(TableCells).BindFunc(func(e *core.RecordEvent) error {
		c.add(e.Record)
		c.sort()
		return e.Next()
	})
	c.app.OnRecordAfterDeleteSuccess(TableCells).BindFunc(func(e *core.RecordEvent) error {
		c.delete(e.Record)
		c.sort()
		return e.Next()
	})
}

func (c *Cells) fetch() error {
	c.cells.Clear()
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
		c.add(cell)
	}
	c.sort()

	return nil
}

func (c *Cells) add(record *core.Record) {
	c.cells.Set(record.Id, record)
	if code := record.GetString("code"); code != "" {
		c.cellsByCode.Set(code, record)
	}
}

func (c *Cells) delete(record *core.Record) {
	c.cells.Delete(record.Id)
	if code := record.GetString("code"); code != "" {
		c.cellsByCode.Delete(code)
	}
}

func (c *Cells) sort() {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.cellsOrder = c.cells.Keys()
	sort.Slice(c.cellsOrder, func(i, j int) bool {
		cell1, _ := c.cells.Get(c.cellsOrder[i])
		cell2, _ := c.cells.Get(c.cellsOrder[j])
		return cell1.GetInt("sort") < cell2.GetInt("sort")
	})

	c.cellsIdsOrder = make(map[string]int, len(c.cellsOrder))
	for key, cellId := range c.cellsOrder {
		c.cellsIdsOrder[cellId] = key
	}
}

// GetByOrder
// Note: cells order starts from 0
func (c *Cells) GetByOrder(order int) (*core.Record, bool) {
	if cellId := c.cellsOrder[order]; cellId != "" {
		return c.cells.Get(cellId)
	}
	return nil, false
}

// GetOrderById
// Note: cells order starts from 0
func (c *Cells) GetOrderById(cellId string) (int, bool) {
	if order, ok := c.cellsIdsOrder[cellId]; ok {
		return order, true
	}
	return 0, false
}

// GetOrderByType
// Note: cells order starts from 0
func (c *Cells) GetOrderByType(t string) (int, bool) {
	if cell, ok := c.GetByType(t); ok {
		return c.GetOrderById(cell.Id)
	}
	return 0, false
}

func (c *Cells) GetByCode(code string) (*core.Record, bool) {
	return c.cellsByCode.Get(code)
}

func (c *Cells) GetAllByType(t string) []*core.Record {
	var res []*core.Record
	for _, record := range c.cells.GetAll() {
		if record.GetString("type") == t {
			res = append(res, record)
		}
	}
	return res
}

func (c *Cells) GetByType(t string) (*core.Record, bool) {
	for _, record := range c.cells.GetAll() {
		if record.GetString("type") == t {
			return record, true
		}
	}
	return nil, false
}

func (c *Cells) Count() int {
	return c.cells.Count()
}

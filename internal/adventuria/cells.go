package adventuria

import (
	"adventuria/pkg/cache"
	"adventuria/pkg/helper"
	"slices"
	"sort"
	"sync"

	"github.com/pocketbase/pocketbase/core"
)

type Cells struct {
	cells         *cache.MemoryCache[string, Cell]
	cellsByCode   *cache.MemoryCache[string, Cell]
	cellsOrder    []string
	cellsIdsOrder map[string]int
	mx            sync.Mutex
}

func NewCells() *Cells {
	cells := &Cells{
		cells:       cache.NewMemoryCache[string, Cell](0, true),
		cellsByCode: cache.NewMemoryCache[string, Cell](0, true),
	}

	cells.fetch()
	cells.bindHooks()

	return cells
}

func (c *Cells) bindHooks() {
	PocketBase.OnRecordAfterCreateSuccess(TableCells).BindFunc(func(e *core.RecordEvent) error {
		err := c.add(e.Record)
		if err != nil {
			return err
		}
		return e.Next()
	})
	PocketBase.OnRecordAfterUpdateSuccess(TableCells).BindFunc(func(e *core.RecordEvent) error {
		err := c.add(e.Record)
		if err != nil {
			return err
		}
		return e.Next()
	})
	PocketBase.OnRecordAfterDeleteSuccess(TableCells).BindFunc(func(e *core.RecordEvent) error {
		c.delete(e.Record)
		return e.Next()
	})
}

func (c *Cells) fetch() error {
	c.cells.Clear()
	c.cellsByCode.Clear()

	cells, err := PocketBase.FindRecordsByFilter(
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

func (c *Cells) add(record *core.Record) error {
	cell, err := NewCellFromRecord(record)
	if err != nil {
		return err
	}

	c.cells.Set(record.Id, cell)
	if code := record.GetString("code"); code != "" {
		c.cellsByCode.Set(code, cell)
	}

	c.sort()

	return nil
}

func (c *Cells) delete(record *core.Record) {
	c.cells.Delete(record.Id)
	if code := record.GetString("code"); code != "" {
		c.cellsByCode.Delete(code)
	}

	c.sort()
}

func (c *Cells) sort() {
	c.mx.Lock()
	defer c.mx.Unlock()

	// Exclude inactive cells from sorting.
	c.cellsOrder = helper.FilterByField(c.cells.Keys(), c.getUnactiveCellsIds(), func(cellId string) string {
		return cellId
	})

	sort.Slice(c.cellsOrder, func(i, j int) bool {
		cell1, _ := c.cells.Get(c.cellsOrder[i])
		cell2, _ := c.cells.Get(c.cellsOrder[j])
		return cell1.Sort() < cell2.Sort()
	})

	c.cellsIdsOrder = make(map[string]int, len(c.cellsOrder))
	for key, cellId := range c.cellsOrder {
		c.cellsIdsOrder[cellId] = key
	}
}

func (c *Cells) getUnactiveCellsIds() []string {
	var res []string
	for _, cell := range c.cells.GetAll() {
		if !cell.IsActive() {
			res = append(res, cell.ID())
		}
	}
	return res
}

// GetByOrder
// Note: cells order starts from 0
func (c *Cells) GetByOrder(order int) (Cell, bool) {
	if order < 0 || order >= len(c.cellsOrder) {
		return nil, false
	}

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
func (c *Cells) GetOrderByType(t CellType) (int, bool) {
	if cell, ok := c.GetByType(t); ok {
		return c.GetOrderById(cell.ID())
	}
	return 0, false
}

func (c *Cells) GetByCode(code string) (Cell, bool) {
	return c.cellsByCode.Get(code)
}

func (c *Cells) GetAllByType(t CellType) []Cell {
	var res []Cell
	for _, cell := range c.cells.GetAll() {
		if cell.IsActive() && cell.Type() == t {
			res = append(res, cell)
		}
	}
	return res
}

func (c *Cells) GetAllByTypes(t []CellType) []Cell {
	var res []Cell
	for _, cell := range c.cells.GetAll() {
		if cell.IsActive() && slices.Contains(t, cell.Type()) {
			res = append(res, cell)
		}
	}
	return res
}

func (c *Cells) GetByType(t CellType) (Cell, bool) {
	for _, cell := range c.cells.GetAll() {
		if cell.IsActive() && cell.Type() == t {
			return cell, true
		}
	}
	return nil, false
}

func (c *Cells) GetById(id string) (Cell, bool) {
	return c.cells.Get(id)
}

func (c *Cells) Count() int {
	return c.cells.Count()
}

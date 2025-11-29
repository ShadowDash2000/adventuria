package adventuria

import (
	"adventuria/pkg/cache"
	"iter"
	"slices"
	"sort"
	"sync"

	"github.com/pocketbase/pocketbase/core"
)

type Cells struct {
	cells         cache.Cache[string, Cell]
	cellsByCode   cache.Cache[string, Cell]
	cellsOrder    []string
	cellsIdsOrder map[string]int
	mx            sync.Mutex
}

func NewCells() (*Cells, error) {
	cells := &Cells{
		cells:       cache.NewMemoryCache[string, Cell](0, true),
		cellsByCode: cache.NewMemoryCache[string, Cell](0, true),
	}

	if err := cells.fetch(); err != nil {
		return nil, err
	}
	cells.bindHooks()

	return cells, nil
}

func (c *Cells) bindHooks() {
	PocketBase.OnRecordAfterCreateSuccess(CollectionCells).BindFunc(func(e *core.RecordEvent) error {
		err := c.add(e.Record)
		if err != nil {
			return err
		}
		return e.Next()
	})
	PocketBase.OnRecordAfterUpdateSuccess(CollectionCells).BindFunc(func(e *core.RecordEvent) error {
		err := c.add(e.Record)
		if err != nil {
			return err
		}
		return e.Next()
	})
	PocketBase.OnRecordAfterDeleteSuccess(CollectionCells).BindFunc(func(e *core.RecordEvent) error {
		c.delete(e.Record)
		return e.Next()
	})
}

func (c *Cells) fetch() error {
	c.cells.Clear()
	c.cellsByCode.Clear()

	cells, err := PocketBase.FindRecordsByFilter(
		CollectionCells,
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

	c.cellsOrder = make([]string, 0, c.cells.Count())
	for _, cell := range c.cells.GetAll() {
		c.cellsOrder = append(c.cellsOrder, cell.ID())
	}

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

// GetOrderByName
// Note: cells order starts from 0
func (c *Cells) GetOrderByName(n string) (int, bool) {
	if cell, ok := c.GetByName(n); ok {
		return c.GetOrderById(cell.ID())
	}
	return 0, false
}

func (c *Cells) GetByCode(code string) (Cell, bool) {
	return c.cellsByCode.Get(code)
}

func (c *Cells) GetAllByType(t CellType) iter.Seq[Cell] {
	return func(yield func(Cell) bool) {
		for _, cell := range c.cells.GetAll() {
			if cell.Type() == t {
				if !yield(cell) {
					return
				}
			}
		}
	}
}

func (c *Cells) GetAllByTypes(t []CellType) iter.Seq[Cell] {
	return func(yield func(Cell) bool) {
		for _, cell := range c.cells.GetAll() {
			if slices.Contains(t, cell.Type()) {
				if !yield(cell) {
					return
				}
			}
		}
	}
}

func (c *Cells) GetByType(t CellType) (Cell, bool) {
	for _, cell := range c.cells.GetAll() {
		if cell.Type() == t {
			return cell, true
		}
	}
	return nil, false
}

func (c *Cells) GetByName(n string) (Cell, bool) {
	for _, cell := range c.cells.GetAll() {
		if cell.Name() == n {
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

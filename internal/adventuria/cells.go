package adventuria

import (
	"adventuria/internal/adventuria/schema"
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

func NewCells(ctx AppContext) (*Cells, error) {
	cells := &Cells{
		cells:       cache.NewMemoryCache[string, Cell](0, true),
		cellsByCode: cache.NewMemoryCache[string, Cell](0, true),
	}

	if err := cells.fetch(ctx); err != nil {
		return nil, err
	}
	cells.bindHooks(ctx)

	return cells, nil
}

func (c *Cells) bindHooks(ctx AppContext) {
	ctx.App.OnRecordAfterCreateSuccess(schema.CollectionCells).BindFunc(func(e *core.RecordEvent) error {
		err := c.add(e.Record)
		if err != nil {
			return err
		}
		return e.Next()
	})
	ctx.App.OnRecordAfterUpdateSuccess(schema.CollectionCells).BindFunc(func(e *core.RecordEvent) error {
		err := c.add(e.Record)
		if err != nil {
			return err
		}
		return e.Next()
	})
	ctx.App.OnRecordAfterDeleteSuccess(schema.CollectionCells).BindFunc(func(e *core.RecordEvent) error {
		c.delete(e.Record)
		return e.Next()
	})
}

func (c *Cells) fetch(ctx AppContext) error {
	c.cells.Clear()
	c.cellsByCode.Clear()

	cells, err := ctx.App.FindRecordsByFilter(
		schema.CollectionCells,
		"",
		"sort",
		0,
		0,
	)
	if err != nil {
		return err
	}

	for _, cell := range cells {
		if err = c.add(cell); err != nil {
			ctx.App.Logger().Error("Cells: unknown cell type", "cell", cell)
		}
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
func (c *Cells) GetOrderByType(t CellType) iter.Seq[int] {
	return func(yield func(int) bool) {
		for cell := range c.GetAllByType(t) {
			order, _ := c.GetOrderById(cell.ID())
			if !yield(order) {
				return
			}
		}
	}
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

func (c *Cells) GetDistanceByIds(firstId, secondId string) (int, bool) {
	var (
		firstCellOrder, secondCellOrder int
		ok                              bool
	)
	if firstCellOrder, ok = c.GetOrderById(firstId); !ok {
		return 0, false
	} else if secondCellOrder, ok = c.GetOrderById(secondId); !ok {
		return 0, false
	}

	return abs(firstCellOrder - secondCellOrder), true
}

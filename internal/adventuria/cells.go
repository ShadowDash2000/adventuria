package adventuria

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/cache"
	"iter"
	"slices"
	"sort"
	"sync"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type Cells struct {
	cells             cache.Cache[string, Cell]
	cellsOrder        map[string][]string
	cellsIdsOrder     map[string]map[string]int
	cellsFlatOrder    []string
	cellsFlatIdsOrder map[string]int
	mx                sync.Mutex
}

func NewCells(ctx AppContext) (*Cells, error) {
	cells := &Cells{
		cells: cache.NewMemoryCache[string, Cell](0, true),
	}

	if err := cells.fetch(ctx); err != nil {
		return nil, err
	}
	cells.bindHooks(ctx)

	return cells, nil
}

func (c *Cells) bindHooks(ctx AppContext) {
	ctx.App.OnRecordAfterCreateSuccess(schema.CollectionCells).BindFunc(func(e *core.RecordEvent) error {
		if disabled := e.Record.GetBool(schema.CellSchema.Disabled); !disabled {
			err := c.add(e.Record)
			if err != nil {
				return err
			}
		}
		return e.Next()
	})
	ctx.App.OnRecordAfterUpdateSuccess(schema.CollectionCells).BindFunc(func(e *core.RecordEvent) error {
		if disabled := e.Record.GetBool(schema.CellSchema.Disabled); !disabled {
			err := c.add(e.Record)
			if err != nil {
				return err
			}
		} else {
			c.delete(e.Record)
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

	var cells []*core.Record
	err := ctx.App.RecordQuery(schema.CollectionCells).
		OrderBy("sort ASC").
		Where(dbx.HashExp{schema.CellSchema.Disabled: false}).
		All(&cells)
	if err != nil {
		return err
	}

	for _, cell := range cells {
		if err = c.addNoSort(cell); err != nil {
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
	c.sort()
	return nil
}

func (c *Cells) addNoSort(record *core.Record) error {
	cell, err := NewCellFromRecord(record)
	if err != nil {
		return err
	}

	c.cells.Set(record.Id, cell)
	return nil
}

func (c *Cells) delete(record *core.Record) {
	c.cells.Delete(record.Id)
	c.sort()
}

func (c *Cells) sort() {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.cellsOrder = make(map[string][]string)
	for _, cell := range c.cells.GetAll() {
		worldId := cell.World()
		c.cellsOrder[worldId] = append(c.cellsOrder[worldId], cell.ID())
	}

	c.cellsIdsOrder = make(map[string]map[string]int)

	for worldId, order := range c.cellsOrder {
		sort.Slice(order, func(i, j int) bool {
			cell1, _ := c.cells.Get(order[i])
			cell2, _ := c.cells.Get(order[j])
			return cell1.Sort() < cell2.Sort()
		})

		c.cellsIdsOrder[worldId] = make(map[string]int, len(order))
		for key, cellId := range order {
			c.cellsIdsOrder[worldId][cellId] = key
		}
	}

	c.cellsFlatOrder = []string{}
	c.cellsFlatIdsOrder = make(map[string]int)

	for _, world := range GameWorlds.GetAll() {
		if order, ok := c.cellsOrder[world.ID()]; ok {
			for _, cellId := range order {
				c.cellsFlatIdsOrder[cellId] = len(c.cellsFlatOrder)
				c.cellsFlatOrder = append(c.cellsFlatOrder, cellId)
			}
		}
	}
}

// GetByOrder
// Note: cells order starts from 0
func (c *Cells) GetByOrder(worldId string, order int) (Cell, bool) {
	orderList, ok := c.cellsOrder[worldId]
	if !ok {
		return nil, false
	}

	if order < 0 || order >= len(orderList) {
		return nil, false
	}

	if cellId := orderList[order]; cellId != "" {
		return c.cells.Get(cellId)
	}
	return nil, false
}

func (c *Cells) GetByGlobalOrder(order int) (Cell, bool) {
	if order < 0 || order >= len(c.cellsFlatOrder) {
		return nil, false
	}

	return c.cells.Get(c.cellsFlatOrder[order])
}

// GetOrderById
// Note: cells order starts from 0
func (c *Cells) GetOrderById(worldId, cellId string) (int, bool) {
	if worldOrders, ok := c.cellsIdsOrder[worldId]; ok {
		if order, ok := worldOrders[cellId]; ok {
			return order, true
		}
	}
	return 0, false
}

func (c *Cells) GetGlobalOrderById(cellId string) (int, bool) {
	order, ok := c.cellsFlatIdsOrder[cellId]
	return order, ok
}

func (c *Cells) Count(worldId string) int {
	return len(c.cellsOrder[worldId])
}

func (c *Cells) CountGlobal() int {
	return len(c.cellsFlatOrder)
}

// GetOrderByType
// Note: cells order starts from 0
func (c *Cells) GetOrderByType(worldId string, t CellType) iter.Seq[int] {
	return func(yield func(int) bool) {
		for cell := range c.GetAllByType(worldId, t) {
			order, _ := c.GetOrderById(worldId, cell.ID())
			if !yield(order) {
				return
			}
		}
	}
}

func (c *Cells) GetGlobalOrderByType(t CellType) iter.Seq[int] {
	return func(yield func(int) bool) {
		for cell := range c.GetAllByTypeGlobal(t) {
			order, _ := c.GetGlobalOrderById(cell.ID())
			if !yield(order) {
				return
			}
		}
	}
}

func (c *Cells) GetAllByType(worldId string, t CellType) iter.Seq[Cell] {
	return func(yield func(Cell) bool) {
		orderList, ok := c.cellsOrder[worldId]
		if !ok {
			return
		}
		for _, cellId := range orderList {
			cell, _ := c.cells.Get(cellId)
			if cell.Type() == t {
				if !yield(cell) {
					return
				}
			}
		}
	}
}

func (c *Cells) GetAllByTypes(worldId string, t []CellType) iter.Seq[Cell] {
	return func(yield func(Cell) bool) {
		orderList, ok := c.cellsOrder[worldId]
		if !ok {
			return
		}
		for _, cellId := range orderList {
			cell, _ := c.cells.Get(cellId)
			if slices.Contains(t, cell.Type()) {
				if !yield(cell) {
					return
				}
			}
		}
	}
}

func (c *Cells) GetAllByTypeGlobal(t CellType) iter.Seq2[Cell, int] {
	return func(yield func(Cell, int) bool) {
		for globalOrder, cellId := range c.cellsFlatOrder {
			cell, _ := c.cells.Get(cellId)
			if cell.Type() == t {
				if !yield(cell, globalOrder) {
					return
				}
			}
		}
	}
}

func (c *Cells) GetById(id string) (Cell, bool) {
	return c.cells.Get(id)
}

func (c *Cells) GetDistanceByIds(worldId, firstId, secondId string) (int, bool) {
	var (
		firstCellOrder, secondCellOrder int
		ok                              bool
	)
	if firstCellOrder, ok = c.GetOrderById(worldId, firstId); !ok {
		return 0, false
	} else if secondCellOrder, ok = c.GetOrderById(worldId, secondId); !ok {
		return 0, false
	}

	return abs(firstCellOrder - secondCellOrder), true
}

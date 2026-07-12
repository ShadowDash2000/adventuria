package cells

import "adventuria/internal/adventuria/model"

type CellDef struct {
	t          model.CellType
	categories []string
	new        CellCreator
}

type CellCreator func(cell model.CellInfo) model.Cell

var registry = &Registry{cellDefs: map[model.CellType]CellDef{}}

type Registry struct {
	cellDefs map[model.CellType]CellDef
}

func (r *Registry) Register(cellDefs ...CellDef) {
	for _, cellDef := range cellDefs {
		r.cellDefs[cellDef.t] = cellDef
	}
}

func (r *Registry) Get(t model.CellType) (CellDef, bool) {
	a, ok := r.cellDefs[t]
	return a, ok
}

func NewCell(t model.CellType, new CellCreator, categories ...string) CellDef {
	return CellDef{
		t:          t,
		categories: categories,
		new:        new,
	}
}

func Register(cellDefs ...CellDef) {
	registry.Register(cellDefs...)
}

func Get(t model.CellType) (CellDef, bool) {
	return registry.Get(t)
}

func Categories(t model.CellType) []string {
	cellDef, ok := registry.Get(t)
	if !ok {
		return nil
	}
	return cellDef.categories
}

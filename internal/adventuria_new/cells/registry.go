package cells

import "adventuria/internal/adventuria_new/model"

type CellDef struct {
	t          model.CellType
	categories []string
	new        CellCreator
}

type CellCreator func(cell model.CellInfo) model.Cell

var registry = &Registry{actions: map[model.CellType]CellDef{}}

type Registry struct {
	actions map[model.CellType]CellDef
}

func (r *Registry) Register(actions ...CellDef) {
	for _, action := range actions {
		r.actions[action.t] = action
	}
}

func (r *Registry) Get(t model.CellType) (CellDef, bool) {
	a, ok := r.actions[t]
	return a, ok
}

func NewCell(t model.CellType, new CellCreator, categories ...string) CellDef {
	return CellDef{
		t:          t,
		categories: categories,
		new:        new,
	}
}

func Register(actions ...CellDef) {
	registry.Register(actions...)
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

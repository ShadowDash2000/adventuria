package actions

import (
	"adventuria/internal/adventuria_new/model"
	"iter"
)

type ActionDef struct {
	t          model.ActionType
	categories []string
	new        ActionCreator
}

func (a *ActionDef) New() model.Action {
	return a.new()
}

func (a *ActionDef) Type() model.ActionType {
	return a.t
}

type ActionCreator func() model.Action

var registry = &Registry{actions: map[model.ActionType]ActionDef{
	//"none": NewAction("none", &NoneAction{}),
}}

type Registry struct {
	actions map[model.ActionType]ActionDef
}

func (r *Registry) Register(actions ...ActionDef) {
	for _, action := range actions {
		r.actions[action.t] = action
	}
}

func (r *Registry) Get(t model.ActionType) (ActionDef, bool) {
	a, ok := r.actions[t]
	return a, ok
}

func (r *Registry) GetAll() iter.Seq2[model.ActionType, ActionDef] {
	return func(yield func(model.ActionType, ActionDef) bool) {
		for t, actionDef := range r.actions {
			if !yield(t, actionDef) {
				return
			}
		}
	}
}

func NewAction(t model.ActionType, new ActionCreator, categories ...string) ActionDef {
	return ActionDef{
		t:          t,
		categories: categories,
		new:        new,
	}
}

func Register(actions ...ActionDef) {
	registry.Register(actions...)
}

func Get(t model.ActionType) (ActionDef, bool) {
	return registry.Get(t)
}

func GetAll() iter.Seq2[model.ActionType, ActionDef] {
	return registry.GetAll()
}

func Categories(t model.ActionType) []string {
	actionDef, ok := registry.Get(t)
	if !ok {
		return nil
	}
	return actionDef.categories
}

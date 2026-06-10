package effects

import (
	"adventuria/internal/adventuria_new/model"
	"fmt"
)

type EffectDef struct {
	t   model.EffectType
	new EffectCreator
}

type EffectCreator func(effect model.EffectInfo) model.Effect

var registry = &Registry{effects: map[model.EffectType]EffectDef{}}

type Registry struct {
	effects map[model.EffectType]EffectDef
}

func (r *Registry) Register(effects ...EffectDef) {
	for _, effect := range effects {
		r.effects[effect.t] = effect
	}
}

func (r *Registry) Get(t model.EffectType) (EffectDef, bool) {
	e, ok := r.effects[t]
	return e, ok
}

func NewEffectDef(t model.EffectType, new EffectCreator) EffectDef {
	return EffectDef{
		t:   t,
		new: new,
	}
}

func Register(effects ...EffectDef) {
	registry.Register(effects...)
}

func Get(t model.EffectType) (EffectDef, bool) {
	return registry.Get(t)
}

func Create(effect model.EffectInfo) (model.Effect, error) {
	effectDef, ok := Get(effect.Type())
	if !ok {
		return nil, fmt.Errorf("effect type %s not registered", effect.Type())
	}
	return effectDef.new(effect), nil
}

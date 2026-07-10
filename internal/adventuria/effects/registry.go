package effects

import (
	"adventuria/internal/adventuria/model"
	"fmt"
	"iter"
)

type EffectDef struct {
	t   model.EffectType
	new EffectCreator
}

type EffectPersistentDef struct {
	t model.EffectType
	e model.EffectPersistent
}

func (e *EffectPersistentDef) Type() model.EffectType {
	return e.t
}

func (e *EffectPersistentDef) Effect() model.EffectPersistent {
	return e.e
}

type EffectCreator func(effect model.EffectInfo) model.Effect

var (
	registry           = &Registry{effects: map[model.EffectType]EffectDef{}}
	registryPersistent = &RegistryPersistent{effects: map[model.EffectType]EffectPersistentDef{}}
)

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

type RegistryPersistent struct {
	effects map[model.EffectType]EffectPersistentDef
}

func (r *RegistryPersistent) Register(effects ...EffectPersistentDef) {
	for _, effect := range effects {
		r.effects[effect.t] = effect
	}
}

func (r *RegistryPersistent) GetAll() iter.Seq2[model.EffectType, EffectPersistentDef] {
	return func(yield func(model.EffectType, EffectPersistentDef) bool) {
		for t, effectDef := range r.effects {
			if !yield(t, effectDef) {
				return
			}
		}
	}
}

func NewEffectDef(t model.EffectType, new EffectCreator) EffectDef {
	return EffectDef{
		t:   t,
		new: new,
	}
}

func NewEffectPersistentDef(t model.EffectType, e model.EffectPersistent) EffectPersistentDef {
	return EffectPersistentDef{
		t: t,
		e: e,
	}
}

func Register(effects ...EffectDef) {
	registry.Register(effects...)
}

func RegisterPersistent(effects ...EffectPersistentDef) {
	registryPersistent.Register(effects...)
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

func GetAllPersistent() iter.Seq2[model.EffectType, EffectPersistentDef] {
	return registryPersistent.GetAll()
}

package adventuria

import (
	"fmt"
	"github.com/pocketbase/pocketbase/core"
	"maps"
	"strings"
)

const (
	EffectUseInstant        = "useInstant"
	EffectUseOnRoll         = "useOnRoll"
	EffectUseOnReroll       = "useOnReroll"
	EffectUseOnDrop         = "useOnDrop"
	EffectUseOnChooseResult = "useOnChooseResult"
	EffectUseOnChooseGame   = "useOnChooseGame"
	EffectUseOnRollItem     = "useOnRollItem"
)

type Effect interface {
	core.RecordProxy
	GetId() string
	Name() string
	Event() string
	Type() string
	Value() any
	AddValue(any)
}

type EffectBase struct {
	core.BaseRecordProxy
	value any
}

var EffectsList = map[string]Effect{}

func RegisterEffects(effects map[string]Effect) {
	maps.Insert(EffectsList, maps.All(effects))
}

func NewEffect(record *core.Record) (Effect, error) {
	t := record.GetString("type")
	fmt.Println("EFFECT:", t, record.Id)
	effect, ok := EffectsList[t]
	if !ok {
		return nil, fmt.Errorf("unknown effect type: %s", t)
	}

	effect.SetProxyRecord(record)

	return effect, nil
}

func (e *EffectBase) GetId() string {
	return e.Id
}

func (e *EffectBase) Name() string {
	return e.GetString("name")
}

func (e *EffectBase) Event() string {
	return e.GetString("event")
}

func (e *EffectBase) Type() string {
	return e.GetString("type")
}

func (e *EffectBase) Value() any {
	return e.value
}

func (e *EffectBase) parseString(s string) []string {
	return strings.Split(s, ", ")
}

type EffectInt struct {
	EffectBase
}

func NewEffectInt() Effect {
	return &EffectInt{
		EffectBase{
			value: 0,
		},
	}
}

func (ei *EffectInt) AddValue(i any) {
	ei.value = ei.value.(int) + i.(int)
}

type EffectBool struct {
	EffectBase
}

func NewEffectBool() Effect {
	return &EffectBool{
		EffectBase{
			value: false,
		},
	}
}

func (eb *EffectBool) AddValue(any) {
	eb.value = true
}

type EffectSlice struct {
	EffectBase
}

func NewEffectSlice() Effect {
	return &EffectSlice{}
}

func (ef *EffectSlice) Value() any {
	if ef.Record == nil {
		return nil
	}

	return ef.parseString(ef.GetString("value"))
}

func (ef *EffectSlice) AddValue(v any) {
	ef.value = v
}

type EffectSliceWithSource[T any] struct {
	EffectBase
	source map[string]T
}

func NewEffectSliceWithSource[T any](source map[string]T) Effect {
	return &EffectSliceWithSource[T]{
		source: source,
	}
}

func (ef *EffectSliceWithSource[T]) Value() any {
	if ef.Record == nil {
		return nil
	}

	var res []T
	sl := ef.parseString(ef.GetString("value"))

	for _, key := range sl {
		if srcVal, ok := ef.source[key]; ok {
			res = append(res, srcVal)
		} else {
			// TODO: log error
		}
	}

	return res
}

func (ef *EffectSliceWithSource[T]) AddValue(v any) {
	ef.value = v
}

type Effects struct {
	effects map[string]Effect
}

func NewEffects() *Effects {
	effects := &Effects{}
	effects.effects = map[string]Effect{}

	for t, effect := range EffectsList {
		effects.effects[t] = effect
	}

	return effects
}

func (ee *Effects) Effect(t string) Effect {
	return ee.effects[t]
}

func (ee *Effects) AddValue(t string, v any) {
	ee.effects[t].AddValue(v)
}

func (ee *Effects) Map() map[string]any {
	m := make(map[string]any)
	for t, effect := range ee.effects {
		m[t] = effect.Value()
	}
	return m
}

func EffectAs[T any](e Effect) T {
	v, _ := e.Value().(T)
	return v
}

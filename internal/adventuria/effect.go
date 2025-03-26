package adventuria

import (
	"fmt"
	"github.com/pocketbase/pocketbase/core"
	"maps"
	"strings"
)

type EffectKind uint

const (
	Int EffectKind = 1 << iota
	Bool
	Slice
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
	Kind() EffectKind
	Type() string
	Value() any
}

type BaseEffect struct {
	core.BaseRecordProxy
	kind EffectKind
}

var EffectsList = map[string]EffectCreator{}

func RegisterEffects(effects map[string]EffectCreator) {
	maps.Insert(EffectsList, maps.All(effects))
}

type EffectCreator func() Effect

func NewEffect(record *core.Record) (Effect, error) {
	effectType := record.GetString("type")
	effectCreator, ok := EffectsList[effectType]
	if !ok {
		return nil, fmt.Errorf("unknown effect type: %s", effectType)
	}

	effect := effectCreator()
	effect.SetProxyRecord(record)

	return effect, nil
}

func (e *BaseEffect) GetId() string {
	return e.Id
}

func (e *BaseEffect) Name() string {
	return e.GetString("name")
}

func (e *BaseEffect) Event() string {
	return e.GetString("event")
}

func (e *BaseEffect) Kind() EffectKind {
	return e.kind
}

func (e *BaseEffect) Type() string {
	return e.GetString("type")
}

func (e *BaseEffect) Value() any {
	return nil
}

func (e *BaseEffect) parseString(s string) []string {
	return strings.Split(s, ", ")
}

type EffectInt struct {
	BaseEffect
}

func NewEffectInt() EffectCreator {
	return func() Effect {
		return &EffectInt{
			BaseEffect: BaseEffect{
				kind: Int,
			},
		}
	}
}

func (ei *EffectInt) Value() any {
	return ei.GetInt("value")
}

type EffectBool struct {
	BaseEffect
}

func NewEffectBool() EffectCreator {
	return func() Effect {
		return &EffectInt{
			BaseEffect: BaseEffect{
				kind: Bool,
			},
		}
	}
}

func (eb *EffectBool) Value() any {
	return eb.GetBool("value")
}

type EffectSlice struct {
	BaseEffect
}

func NewEffectSlice() EffectCreator {
	return func() Effect {
		return &EffectSlice{
			BaseEffect: BaseEffect{
				kind: Slice,
			},
		}
	}
}

func (ef *EffectSlice) Value() any {
	return ef.parseString(ef.GetString("value"))
}

type EffectSliceWithSource[T any] struct {
	BaseEffect
	source map[string]T
}

func NewEffectSliceWithSource[T any](source map[string]T) EffectCreator {
	return func() Effect {
		return &EffectSliceWithSource[T]{
			BaseEffect: BaseEffect{
				kind: Slice,
			},
			source: source,
		}
	}
}

func (ef *EffectSliceWithSource[T]) Value() any {
	var res []any
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

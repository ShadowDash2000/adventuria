package adventuria

import (
	"fmt"
	"github.com/pocketbase/pocketbase/core"
	"maps"
)

type Kind uint

const (
	Int Kind = iota
	Bool
	String
	Slice
	SliceWithSource
)

type EffectUse string

const (
	EffectUseInstant        EffectUse = "useInstant"
	EffectUseOnRoll         EffectUse = "useOnRoll"
	EffectUseOnReroll       EffectUse = "useOnReroll"
	EffectUseOnDrop         EffectUse = "useOnDrop"
	EffectUseOnChooseResult EffectUse = "useOnChooseResult"
	EffectUseOnChooseGame   EffectUse = "useOnChooseGame"
	EffectUseOnRollItem     EffectUse = "useOnRollItem"
	EffectUseOnAny          EffectUse = "useOnAny"
)

type Effect interface {
	core.RecordProxy
	ID() string
	Name() string
	Event() EffectUse
	Type() string
	Condition() string
	Int() int
	Bool() bool
	String() string
	Slice() []string
	ParseJSON()
	AddInt(int)
	AddBool(bool)
	AddString(string)
	AddSlice([]string)
	Kind() Kind
	ReceiveSource() any
}

type EffectBase struct {
	core.BaseRecordProxy
	value  EffectValue
	kind   Kind
	source EffectSourceReceiver
	event  EffectUse
}

type EffectValue struct {
	Int    int      `json:"int"`
	Bool   bool     `json:"bool"`
	String string   `json:"string"`
	Slice  []string `json:"slice"`
}

var EffectsList = map[string]EffectCreator{}

type EffectCreator func() Effect

func RegisterEffects(effects map[string]EffectCreator) {
	maps.Insert(EffectsList, maps.All(effects))
}

func NewEffectRecord(record *core.Record) (Effect, error) {
	t := record.GetString("type")

	effectCreator, ok := EffectsList[t]
	if !ok {
		return nil, fmt.Errorf("unknown effect type: %s", t)
	}

	effect := effectCreator()
	effect.SetProxyRecord(record)
	effect.ParseJSON()

	return effect, nil
}

func NewEffect(kind Kind, event EffectUse) EffectCreator {
	return func() Effect {
		return &EffectBase{
			value: EffectValue{},
			kind:  kind,
			event: event,
		}
	}
}

func (e *EffectBase) ID() string {
	return e.Id
}

func (e *EffectBase) Name() string {
	return e.GetString("name")
}

func (e *EffectBase) Event() EffectUse {
	if e.event != EffectUseOnAny {
		return e.event
	}

	if event := e.GetString("useOn"); event != "" {
		return EffectUse(event)
	}

	return e.event
}

func (e *EffectBase) Type() string {
	return e.GetString("type")
}

func (e *EffectBase) Condition() string {
	return e.GetString("condition")
}

func (e *EffectBase) Int() int {
	return e.value.Int
}

func (e *EffectBase) Bool() bool {
	return e.value.Bool
}

func (e *EffectBase) String() string {
	return e.value.String
}

func (e *EffectBase) Slice() []string {
	return e.value.Slice
}

func (e *EffectBase) ParseJSON() {
	e.UnmarshalJSONField("value", &e.value)
}

func (e *EffectBase) AddInt(i int) {
	e.value.Int += i
}

func (e *EffectBase) AddBool(b bool) {
	e.value.Bool = b
}

func (e *EffectBase) AddString(s string) {
	e.value.String = s
}

func (e *EffectBase) AddSlice(slice []string) {
	e.value.Slice = slice
}

func (e *EffectBase) Kind() Kind {
	return e.kind
}

func (e *EffectBase) ReceiveSource() any {
	return e.source(e.Slice())
}

type EffectSourceGiver[T any] interface {
	Slice() []T
}

type EffectSourceReceiver func([]string) any

func DefaultEffectSourceReceiver(source []string) any {
	return source
}

func NewEffectWithSource(source EffectSourceReceiver, event EffectUse) EffectCreator {
	return func() Effect {
		return &EffectBase{
			value:  EffectValue{},
			kind:   SliceWithSource,
			source: source,
			event:  event,
		}
	}
}

type Effects struct {
	effects map[string]Effect
}

// NewEffects
// My dear diary, I am losing my mind with this piece of fucking
// retarded thrash of bullshit. Sorry, I'm just too tired right now.
// It's 4AM... Maybe I should go to sleep now, but I'm doing this stupid Effects.
func NewEffects() *Effects {
	effects := &Effects{}

	effects.effects = make(map[string]Effect, len(EffectsList))
	for t, creator := range EffectsList {
		effects.effects[t] = creator()
	}

	return effects
}

func (ee *Effects) Add(effect Effect) {
	e, ok := ee.effects[effect.Type()]
	if !ok {
		return
	}

	switch e.Kind() {
	case Int:
		e.AddInt(effect.Int())
	case Bool:
		e.AddBool(true)
	case String:
		e.AddString(effect.String())
	case Slice, SliceWithSource:
		e.AddSlice(effect.Slice())
	}

	ee.effects[effect.Type()] = e
}

func (ee *Effects) Effect(t string) Effect {
	return ee.effects[t]
}

func (ee *Effects) Map() map[string]any {
	m := make(map[string]any)
	for t, e := range ee.effects {
		switch e.Kind() {
		case Int:
			m[t] = e.Int()
		case Bool:
			m[t] = e.Bool()
		case String:
			m[t] = e.String()
		case Slice:
			m[t] = e.Slice()
		case SliceWithSource:
			m[t] = e.ReceiveSource()
		}
	}
	return m
}

package adventuria

import (
	"fmt"
	"github.com/pocketbase/pocketbase/core"
	"strings"
)

type EffectKind uint

const (
	Int EffectKind = 1 << iota
	Bool
	Slice
)

const (
	ItemUseInstant        = "useInstant"
	ItemUseOnRoll         = "useOnRoll"
	ItemUseOnReroll       = "useOnReroll"
	ItemUseOnDrop         = "useOnDrop"
	ItemUseOnChooseResult = "useOnChooseResult"
	ItemUseOnChooseGame   = "useOnChooseGame"
	ItemUseOnRollItem     = "useOnRollItem"

	EffectTypeNothing          = "nothing"
	EffectTypePointsIncrement  = "pointsIncrement"
	EffectTypeJailEscape       = "jailEscape"
	EffectTypeDiceMultiplier   = "diceMultiplier"
	EffectTypeDiceIncrement    = "diceIncrement"
	EffectTypeChangeDices      = "changeDices"
	EffectTypeSafeDrop         = "isSafeDrop"
	EffectTypeTimerIncrement   = "timerIncrement"
	EffectTypeRollReverse      = "rollReverse"
	EffectTypeDropInventory    = "dropInventory"
	EffectTypeCellPointsDivide = "cellPointsDivide"
)

var (
	EffectKindInt IEffect = &EffectInt{Effect{
		kind: Int,
	}}
	EffectKindBool IEffect = &EffectBool{Effect{
		kind: Bool,
	}}
	EffectKindSlice IEffect = &EffectSlice{Effect{
		kind: Slice,
	}}
)

var EffectsKindList = map[string]IEffect{
	EffectTypeNothing:          EffectKindInt,
	EffectTypePointsIncrement:  EffectKindInt,
	EffectTypeJailEscape:       EffectKindBool,
	EffectTypeDiceMultiplier:   EffectKindInt,
	EffectTypeDiceIncrement:    EffectKindInt,
	EffectTypeChangeDices:      EffectKindSlice,
	EffectTypeSafeDrop:         EffectKindBool,
	EffectTypeTimerIncrement:   EffectKindInt,
	EffectTypeRollReverse:      EffectKindBool,
	EffectTypeDropInventory:    EffectKindBool,
	EffectTypeCellPointsDivide: EffectKindInt,
}

type IEffect interface {
	SetRecord(*core.Record)
	Id() string
	Name() string
	Event() string
	Kind() EffectKind
	Type() string
	Value() any
}

type Effect struct {
	effect *core.Record
	kind   EffectKind
}

func NewEffect(record *core.Record) (IEffect, error) {
	effectType := record.GetString("type")
	effect, ok := EffectsKindList[effectType]
	if !ok {
		return nil, fmt.Errorf("unknown effect type: %s", effectType)
	}

	effect.SetRecord(record)

	return effect, nil
}

func (e *Effect) SetRecord(record *core.Record) {
	e.effect = record
}

func (e *Effect) Id() string {
	return e.effect.Id
}

func (e *Effect) Name() string {
	return e.effect.GetString("name")
}

func (e *Effect) Event() string {
	return e.effect.GetString("event")
}

func (e *Effect) Kind() EffectKind {
	return e.kind
}

func (e *Effect) Type() string {
	return e.effect.GetString("type")
}

func (e *Effect) Value() any {
	return nil
}

func (e *Effect) parseString(s string) []string {
	return strings.Split(s, ", ")
}

type EffectInt struct {
	Effect
}

func (ei *EffectInt) Value() any {
	return ei.effect.GetInt("value")
}

type EffectBool struct {
	Effect
}

func (eb *EffectBool) Value() any {
	return eb.effect.GetBool("value")
}

type EffectSlice struct {
	Effect
}

func (ef *EffectSlice) Value() any {
	var res []any
	sl := ef.parseString(ef.effect.GetString("value"))

	// TODO: instead of switch case here, this should be different Effect struct for different effects of type slice
	// Or maybe, we can somehow pass a slice from which we can append elements
	switch ef.Type() {
	case EffectTypeChangeDices:
		for _, v := range sl {
			res = append(res, Dices[v])
		}
	}

	return res
}

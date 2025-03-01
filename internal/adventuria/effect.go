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
	ItemUseOnRollItem     = "useOnRollItem"

	EffectTypeNothing         = "nothing"
	EffectTypePointsIncrement = "pointsIncrement"
	EffectTypeJailEscape      = "jailEscape"
	EffectTypeDiceMultiplier  = "diceMultiplier"
	EffectTypeDiceIncrement   = "diceIncrement"
	EffectTypeChangeDices     = "changeDices"
	EffectTypeSafeDrop        = "safeDrop"
	EffectTypeTimerIncrement  = "timerIncrement"
	EffectTypeRollReverse     = "rollReverse"
	EffectTypeDropInventory   = "dropInventory"
)

var (
	EffectsKindList = map[string]EffectKind{
		EffectTypePointsIncrement: Int,
		EffectTypeJailEscape:      Bool,
		EffectTypeDiceMultiplier:  Int,
		EffectTypeDiceIncrement:   Int,
		EffectTypeChangeDices:     Slice,
		EffectTypeSafeDrop:        Bool,
		EffectTypeTimerIncrement:  Int,
		EffectTypeRollReverse:     Bool,
	}

	InstantsEffectsList = map[string]struct{}{
		EffectTypePointsIncrement: {},
		EffectTypeJailEscape:      {},
		EffectTypeTimerIncrement:  {},
	}
	OnRollEffectsList = map[string]struct{}{
		EffectTypeDiceMultiplier: {},
		EffectTypeDiceIncrement:  {},
		EffectTypeChangeDices:    {},
		EffectTypeRollReverse:    {},
	}
	OnRerollEffectsList = map[string]struct{}{}
	OnDropEffectsList   = map[string]struct{}{
		EffectTypeSafeDrop: {},
	}
	OnChooseResultEffectsList = map[string]struct{}{
		EffectTypeNothing: {},
	}
	OnRollItemEffectsList = map[string]struct{}{
		EffectTypeDropInventory: {},
	}

	EffectsList = map[string]map[string]struct{}{
		ItemUseInstant:        InstantsEffectsList,
		ItemUseOnRoll:         OnRollEffectsList,
		ItemUseOnReroll:       OnRerollEffectsList,
		ItemUseOnDrop:         OnDropEffectsList,
		ItemUseOnChooseResult: OnChooseResultEffectsList,
		ItemUseOnRollItem:     OnRollItemEffectsList,
	}
)

type Effect struct {
	effect *core.Record
	kind   EffectKind
}

func NewEffect(record *core.Record) (*Effect, error) {
	effect := &Effect{
		effect: record,
	}

	var ok bool
	if effect.kind, ok = EffectsKindList[effect.Type()]; !ok {
		return nil, fmt.Errorf("unknown effect type: %s", effect.Type())
	}

	return effect, nil
}

func (e *Effect) Id() string {
	return e.effect.Id
}

func (e *Effect) Kind() EffectKind {
	return e.kind
}

func (e *Effect) Type() string {
	return e.effect.GetString("type")
}

func (e *Effect) GetInt() int {
	return e.effect.GetInt("value")
}

func (e *Effect) GetSlice() []any {
	var res []any
	sl := e.parseString(e.effect.GetString("value"))

	switch e.Type() {
	case EffectTypeChangeDices:
		for _, v := range sl {
			res = append(res, Dices[v])
		}
	}

	return res
}

func (e *Effect) parseString(s string) []string {
	return strings.Split(s, ", ")
}

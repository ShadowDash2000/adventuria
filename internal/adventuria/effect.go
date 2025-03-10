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
	EffectsKindList = map[string]EffectKind{
		EffectTypeNothing:          Int,
		EffectTypePointsIncrement:  Int,
		EffectTypeJailEscape:       Bool,
		EffectTypeDiceMultiplier:   Int,
		EffectTypeDiceIncrement:    Int,
		EffectTypeChangeDices:      Slice,
		EffectTypeSafeDrop:         Bool,
		EffectTypeTimerIncrement:   Int,
		EffectTypeRollReverse:      Bool,
		EffectTypeDropInventory:    Bool,
		EffectTypeCellPointsDivide: Int,
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

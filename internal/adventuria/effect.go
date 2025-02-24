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
	EffectTypeDiceMultiplier = "diceMultiplier"
	EffectTypeDiceIncrement  = "diceIncrement"
	EffectTypeSafeDrop       = "safeDrop"
	EffectTypeChangeDices    = "changeDices"
)

var EffectsList = map[string]EffectKind{
	EffectTypeDiceMultiplier: Int,
	EffectTypeDiceIncrement:  Int,
	EffectTypeSafeDrop:       Bool,
	EffectTypeChangeDices:    Slice,
}

type Effect struct {
	effect *core.Record
	kind   EffectKind
}

func NewEffect(record *core.Record) (*Effect, error) {
	effect := &Effect{
		effect: record,
	}

	var ok bool
	if effect.kind, ok = EffectsList[effect.Type()]; !ok {
		return nil, fmt.Errorf("unknown effect type: %s", effect.Type())
	}

	return effect, nil
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

func (e *Effect) GetSlice() []string {
	return e.parseString(e.effect.GetString("value"))
}

func (e *Effect) parseString(s string) []string {
	return strings.Split(s, ", ")
}

package adventuria

import (
	"math/rand/v2"
)

type Dice interface {
	Roll() int
	Type() string
	MaxNumber() int
}

type DiceResponse struct {
	Type      string `json:"type"`
	MaxNumber int    `json:"maxNumber"`
}

type DiceBase struct {
	t         string
	maxNumber int
}

func (d *DiceBase) Roll() int {
	return rand.IntN(d.maxNumber) + 1
}

func (d *DiceBase) Type() string {
	return d.t
}

func (d *DiceBase) MaxNumber() int {
	return d.maxNumber
}

var DicesList = map[string]Dice{
	"d4": DiceTypeD4,
	"d6": DiceTypeD6,
	"d8": DiceTypeD8,
}

var (
	DiceTypeD4 = &DiceBase{"d4", 4}
	DiceTypeD6 = &DiceBase{"d6", 6}
	DiceTypeD8 = &DiceBase{"d8", 8}
)

type DiceEffectSourceGiver[T Dice] struct {
	source []string
}

func NewDiceEffectSourceGiver[T Dice](source []string) EffectSourceGiver[T] {
	return &DiceEffectSourceGiver[T]{source: source}
}

func (dg *DiceEffectSourceGiver[T]) Slice() []T {
	var res []T
	for _, key := range dg.source {
		if dice, ok := DicesList[key]; ok {
			res = append(res, dice.(T))
		} else {
			// TODO: log error
		}
	}
	return res
}

func DiceEffectSourceReceiver(source []string) any {
	var res []any
	for _, key := range source {
		if dice, ok := DicesList[key]; ok {
			res = append(res, &DiceResponse{
				Type:      dice.Type(),
				MaxNumber: dice.MaxNumber(),
			})
		} else {
			// TODO: log error
		}
	}
	return res
}

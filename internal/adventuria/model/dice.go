package model

import "math/rand/v2"

type DiceType string

const (
	DiceTypeD4 DiceType = "d4"
	DiceTypeD6 DiceType = "d6"
)

type Dice int

func (d Dice) Roll() int {
	if d <= 0 {
		return 0
	}
	return rand.IntN(int(d)) + 1
}

func GetDice(diceType DiceType) (Dice, bool) {
	dice, ok := diceList[diceType]
	return dice, ok
}

var diceList = map[DiceType]Dice{
	DiceTypeD4: diceD4,
	DiceTypeD6: diceD6,
}

var (
	diceD4 = Dice(4)
	diceD6 = Dice(6)
)

func DiceD4() Dice {
	return diceD4
}

func DiceD6() Dice {
	return diceD6
}

package adventuria

import "math/rand/v2"

type Dice struct {
	Type      string `json:"type"`
	MaxNumber int    `json:"maxNumber"`
}

var Dices = map[string]Dice{
	"d4": DiceTypeD4,
	"d6": DiceTypeD6,
	"d8": DiceTypeD8,
}

var (
	DiceTypeD4 Dice = Dice{"d4", 4}
	DiceTypeD6 Dice = Dice{"d6", 6}
	DiceTypeD8 Dice = Dice{"d8", 8}
)

func (d *Dice) Roll() int {
	return rand.IntN(d.MaxNumber-1) + 1
}

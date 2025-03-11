package adventuria

type DropEffects struct {
	IsSafeDrop bool
}

type DoneEffects struct {
	CellPointsDivide int
}

type RollEffects struct {
	DiceMultiplier int
	DiceIncrement  int
	Dices          []Dice
	RollReverse    bool
}

type RollResult struct {
	n int
}

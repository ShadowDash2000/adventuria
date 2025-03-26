package adventuria

type DropEffects struct {
	IsSafeDrop bool
}

type DoneEffects struct {
	CellPointsDivide int
}

type RollDicesResult struct {
	Dices []Dice
}

type RollResult struct {
	N int
}

package adventuria

const (
	TableUsers     = "users"
	TableActions   = "actions"
	TableCells     = "cells"
	TableItems     = "items"
	TableInventory = "inventory"

	ActionTypeRoll       = "roll"
	ActionTypeReroll     = "reroll"
	ActionTypeDrop       = "drop"
	ActionTypeDone       = "done"
	ActionTypeGame       = "game"
	ActionTypeRollCell   = "rollCell"
	ActionTypeRollPreset = "rollPreset"

	CellTypeGame   = "game"
	CellTypeStart  = "start"
	CellTypeJail   = "jail"
	CellTypeBigWin = "big-win"
	CellTypePreset = "preset"

	UserNextStepRoll             = "roll"
	UserNextStepChooseResult     = "chooseResult"
	UserNextStepChooseGame       = "chooseGame"
	UserNextStepRollJailCell     = "rollJailCell"
	UserNextStepRollBigWin       = "rollBigWin"
	UserNextStepRollOnBigWinDrop = "rollOnBigWinDrop"
	UserNextStepRollPreset       = "rollPreset"
)

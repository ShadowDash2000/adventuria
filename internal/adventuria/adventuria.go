package adventuria

const (
	TableUsers      = "users"
	TableActions    = "actions"
	TableCells      = "cells"
	TableItems      = "items"
	TableInventory  = "inventory"
	TableWheelItems = "wheel_items"

	ActionTypeRoll       = "roll"
	ActionTypeReroll     = "reroll"
	ActionTypeDrop       = "drop"
	ActionTypeDone       = "done"
	ActionTypeGame       = "game"
	ActionTypeRollCell   = "rollCell"
	ActionTypeRollPreset = "rollPreset"
	ActionTypeRollMovie  = "rollMovie"

	CellTypeGame   = "game"
	CellTypeStart  = "start"
	CellTypeJail   = "jail"
	CellTypeBigWin = "big-win"
	CellTypePreset = "preset"
	CellTypeMovie  = "movie"

	UserNextStepRoll             = "roll"
	UserNextStepChooseResult     = "chooseResult"
	UserNextStepChooseGame       = "chooseGame"
	UserNextStepRollJailCell     = "rollJailCell"
	UserNextStepRollBigWin       = "rollBigWin"
	UserNextStepRollOnBigWinDrop = "rollOnBigWinDrop"
	UserNextStepRollPreset       = "rollPreset"
	UserNextStepRollMovie        = "rollMovie"
	UserNextStepMovieResult      = "movieResult"
)

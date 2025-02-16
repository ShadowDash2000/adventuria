package adventuria

const (
	TableUsers          = "users"
	TableActions        = "actions"
	TableCells          = "cells"
	TableItems          = "items"
	TableInventory      = "inventory"
	TableActionsEffects = "actionsEffects"

	ActionStatusNone          = "none"
	ActionStatusGameNotChosen = "gameNotChosen"
	ActionStatusReroll        = "reroll"
	ActionStatusDrop          = "drop"
	ActionStatusDone          = "done"
	ActionStatusInProgress    = "inProgress"

	CellTypeGame   = "game"
	CellTypeStart  = "start"
	CellTypeJail   = "jail"
	CellTypeBigWin = "big-win"
	CellTypePreset = "preset"

	UserNextStepRoll         = "roll"
	UserNextStepChooseResult = "chooseResult"
	UserNextStepChooseGame   = "chooseGame"
	UserNextStepRollCell     = "rollCell"
)

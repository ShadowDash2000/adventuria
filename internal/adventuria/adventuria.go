package adventuria

const (
	TableUsers      = "users"
	TableActions    = "actions"
	TableCells      = "cells"
	TableItems      = "items"
	TableInventory  = "inventory"
	TableWheelItems = "wheel_items"
	TableLogs       = "logs"
	TableTimers     = "timers"
	TableSettings   = "settings"

	ActionTypeRoll            = "roll"
	ActionTypeReroll          = "reroll"
	ActionTypeDrop            = "drop"
	ActionTypeChooseResult    = "chooseResult"
	ActionTypeChooseGame      = "chooseGame"
	ActionTypeRollCell        = "rollCell"
	ActionTypeRollItem        = "rollItem"
	ActionTypeRollWheelPreset = "rollWheelPreset"

	CellTypeGame        = "game"
	CellTypeStart       = "start"
	CellTypeJail        = "jail"
	CellTypePreset      = "preset"
	CellTypeItem        = "item"
	CellTypeWheelPreset = "wheelPreset"

	LogTypeItemUse           = "itemUse"
	LogTypeItemDrop          = "itemDrop"
	LogTypeItemEffectApplied = "itemEffectApplied"
)

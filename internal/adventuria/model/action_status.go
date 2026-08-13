package model

type ActionStatus string

const (
	ActionStatusNone            ActionStatus = "none"
	ActionStatusStart           ActionStatus = "start"
	ActionStatusMove            ActionStatus = "move"
	ActionStatusRollDice        ActionStatus = "roll_dice"
	ActionStatusDone            ActionStatus = "done"
	ActionStatusReroll          ActionStatus = "reroll"
	ActionStatusDrop            ActionStatus = "drop"
	ActionStatusRollWheel       ActionStatus = "roll_wheel"
	ActionStatusRollItemOnCell  ActionStatus = "roll_item_on_cell"
	ActionStatusTeleport        ActionStatus = "teleport"
	ActionStatusNeedToRollWheel ActionStatus = "need_to_roll_wheel"
)

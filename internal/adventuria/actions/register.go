package actions

import "adventuria/internal/adventuria"

const (
	ActionTypeRollDice       adventuria.ActionType = "rollDice"
	ActionTypeDone           adventuria.ActionType = "done"
	ActionTypeReroll         adventuria.ActionType = "reroll"
	ActionTypeDrop           adventuria.ActionType = "drop"
	ActionTypeRollWheel      adventuria.ActionType = "rollWheel"
	ActionTypeRollItem       adventuria.ActionType = "rollItem"
	ActionTypeBuyItem        adventuria.ActionType = "buyItem"
	ActionTypeUpdateComment  adventuria.ActionType = "update_comment"
	ActionTypeRollItemOnCell adventuria.ActionType = "rollItemOnCell"
)

func WithBaseActions() {
	adventuria.RegisterActions([]adventuria.ActionDef{
		adventuria.NewAction("rollDice", &RollDiceAction{}),
		adventuria.NewAction("done", &DoneAction{}),
		adventuria.NewAction("reroll", &RerollAction{}),
		adventuria.NewAction("drop", &DropAction{}),
		adventuria.NewAction("rollWheel", &RollWheelAction{}, "wheel_roll", "on_cell"),
		adventuria.NewAction("rollItem", &RollItemAction{}, "wheel_roll"),
		adventuria.NewAction("buyItem", &BuyAction{}),
		adventuria.NewAction("update_comment", &UpdateCommentAction{}),
		adventuria.NewAction("rollItemOnCell", &RollItemOnCellAction{}, "wheel_roll", "on_cell"),
		adventuria.NewAction("refreshShop", &RefreshShopAction{}),
	})
}

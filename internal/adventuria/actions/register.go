package actions

import "adventuria/internal/adventuria"

const (
	ActionTypeRollDice  adventuria.ActionType = "rollDice"
	ActionTypeDone      adventuria.ActionType = "done"
	ActionTypeReroll    adventuria.ActionType = "reroll"
	ActionTypeDrop      adventuria.ActionType = "drop"
	ActionTypeRollWheel adventuria.ActionType = "rollWheel"
	ActionTypeRollItem  adventuria.ActionType = "rollItem"
)

func WithBaseActions() {
	adventuria.RegisterActions([]adventuria.ActionCreator{
		adventuria.NewAction(ActionTypeRollDice, &RollDiceAction{}),
		adventuria.NewAction(ActionTypeDone, &DoneAction{}),
		adventuria.NewAction(ActionTypeReroll, &RerollAction{}),
		adventuria.NewAction(ActionTypeDrop, &DropAction{}),
		adventuria.NewAction(ActionTypeRollWheel, &RollWheelAction{}),
		adventuria.NewAction(ActionTypeRollItem, &RollItemAction{}),
	})
}

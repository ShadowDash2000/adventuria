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
	adventuria.RegisterActions(map[adventuria.ActionType]adventuria.ActionCreator{
		ActionTypeRollDice:  adventuria.NewAction(&RollDiceAction{}),
		ActionTypeDone:      adventuria.NewAction(&DoneAction{}),
		ActionTypeReroll:    adventuria.NewAction(&RerollAction{}),
		ActionTypeDrop:      adventuria.NewAction(&DropAction{}),
		ActionTypeRollWheel: adventuria.NewAction(&RollWheelAction{}),
		ActionTypeRollItem:  adventuria.NewAction(&RollItemAction{}),
	})
}

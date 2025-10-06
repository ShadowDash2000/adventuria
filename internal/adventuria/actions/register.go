package actions

import "adventuria/internal/adventuria"

func WithBaseActions() {
	adventuria.RegisterActions(map[string]adventuria.ActionCreator{
		"rollDice":  adventuria.NewAction(&RollDiceAction{}),
		"done":      adventuria.NewAction(&DoneAction{}),
		"reroll":    adventuria.NewAction(&RerollAction{}),
		"drop":      adventuria.NewAction(&DropAction{}),
		"rollWheel": adventuria.NewAction(&RollWheelAction{}),
		"rollItem":  adventuria.NewAction(&RollItemAction{}),
	})
}

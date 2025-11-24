package actions

import (
	"adventuria/internal/adventuria"
	"fmt"
)

type RollDiceAction struct {
	adventuria.ActionBase
}

func (a *RollDiceAction) CanDo(user adventuria.User) bool {
	return user.LastAction().CanMove()
}

type RollDiceResult struct {
	Roll        int             `json:"roll"`
	DiceRolls   []DiceRoll      `json:"dice_rolls"`
	CurrentCell adventuria.Cell `json:"current_cell"`
}

type DiceRoll struct {
	Type string `json:"type"`
	Roll int    `json:"roll"`
}

func (a *RollDiceAction) Do(user adventuria.User, _ adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	onBeforeRollEvent := &adventuria.OnBeforeRollEvent{
		Dices: []*adventuria.Dice{adventuria.DiceTypeD6, adventuria.DiceTypeD6},
	}
	err := user.OnBeforeRoll().Trigger(onBeforeRollEvent)
	if err != nil {
		adventuria.PocketBase.Logger().Error(
			"roll_dice.do(): failed to trigger onBeforeRoll event",
			"error",
			err,
		)
	}

	onBeforeRollMoveEvent := &adventuria.OnBeforeRollMoveEvent{
		N: 0,
	}
	diceRolls := make([]DiceRoll, len(onBeforeRollEvent.Dices))
	for i, dice := range onBeforeRollEvent.Dices {
		diceRolls[i] = DiceRoll{
			Type: dice.Type,
			Roll: dice.Roll(),
		}
		onBeforeRollMoveEvent.N += diceRolls[i].Roll
	}

	err = user.OnBeforeRollMove().Trigger(onBeforeRollMoveEvent)
	if err != nil {
		adventuria.PocketBase.Logger().Error(
			"roll_dice.do(): failed to trigger onBeforeRollMove event",
			"error",
			err,
		)
	}

	onAfterMoveEvent, err := user.Move(onBeforeRollMoveEvent.N)
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error",
		}, fmt.Errorf("roll_dice.do(): %w", err)
	}

	onAfterRollEvent := &adventuria.OnAfterRollEvent{
		Dices: onBeforeRollEvent.Dices,
		N:     onBeforeRollMoveEvent.N,
	}
	err = user.OnAfterRoll().Trigger(onAfterRollEvent)
	if err != nil {
		adventuria.PocketBase.Logger().Error(
			"roll_dice.do(): failed to trigger onAfterRoll event",
			"error",
			err,
		)
	}

	return &adventuria.ActionResult{
		Success: true,
		Data: RollDiceResult{
			Roll:        onBeforeRollMoveEvent.N,
			DiceRolls:   diceRolls,
			CurrentCell: onAfterMoveEvent.CurrentCell,
		},
	}, nil
}

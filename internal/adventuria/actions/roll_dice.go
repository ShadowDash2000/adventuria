package actions

import (
	"adventuria/internal/adventuria"
)

type RollDiceAction struct {
	adventuria.ActionBase
}

func (a *RollDiceAction) CanDo(user adventuria.User) bool {
	return user.LastAction().CanMove()
}

type RollDiceResult struct {
	Roll        int             `json:"roll"`
	DiceRolls   []int           `json:"dice_roll"`
	CurrentCell adventuria.Cell `json:"current_cell"`
}

func (a *RollDiceAction) Do(user adventuria.User, _ adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	onBeforeRollEvent := &adventuria.OnBeforeRollEvent{
		Dices: []*adventuria.Dice{adventuria.DiceTypeD6, adventuria.DiceTypeD6},
	}
	err := user.OnBeforeRoll().Trigger(onBeforeRollEvent)
	if err != nil {
		return nil, err
	}

	onBeforeRollMoveEvent := &adventuria.OnBeforeRollMoveEvent{
		N: 0,
	}
	diceRolls := make([]int, len(onBeforeRollEvent.Dices))
	for i, dice := range onBeforeRollEvent.Dices {
		diceRolls[i] = dice.Roll()
		onBeforeRollMoveEvent.N += diceRolls[i]
	}

	err = user.OnBeforeRollMove().Trigger(onBeforeRollMoveEvent)
	if err != nil {
		return nil, err
	}

	onAfterMoveEvent, err := user.Move(onBeforeRollMoveEvent.N)
	if err != nil {
		return nil, err
	}

	onAfterRollEvent := &adventuria.OnAfterRollEvent{
		Dices: onBeforeRollEvent.Dices,
		N:     onBeforeRollMoveEvent.N,
	}
	err = user.OnAfterRoll().Trigger(onAfterRollEvent)
	if err != nil {
		return nil, err
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

package actions

import (
	"adventuria/internal/adventuria"
)

type RollDiceAction struct {
	adventuria.ActionBase
}

func (a *RollDiceAction) CanDo() bool {
	return a.User().GetNextStepType() == adventuria.ActionTypeRollDice
}

func (a *RollDiceAction) Do(_ adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	onBeforeRollEvent := &adventuria.OnBeforeRollEvent{
		Dices: []*adventuria.Dice{adventuria.DiceTypeD6, adventuria.DiceTypeD6},
	}
	err := a.User().OnBeforeRoll().Trigger(onBeforeRollEvent)
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

	err = a.User().OnBeforeRollMove().Trigger(onBeforeRollMoveEvent)
	if err != nil {
		return nil, err
	}

	onAfterMoveEvent, err := a.User().Move(onBeforeRollMoveEvent.N)

	onAfterRollEvent := &adventuria.OnAfterRollEvent{
		Dices: onBeforeRollEvent.Dices,
		N:     onBeforeRollMoveEvent.N,
	}
	err = a.User().OnAfterRoll().Trigger(onAfterRollEvent)
	if err != nil {
		return nil, err
	}

	return &adventuria.ActionResult{
		Success: true,
		Data: map[string]interface{}{
			"roll":         onBeforeRollMoveEvent.N,
			"dice_roll":    diceRolls,
			"current_cell": onAfterMoveEvent.CurrentCell,
		},
	}, nil
}

package actions

import (
	"adventuria/internal/adventuria"
	"math/rand/v2"
)

type RollDiceAction struct {
	adventuria.ActionBase
}

func (a *RollDiceAction) CanDo() bool {
	switch a.User().LastAction().Type() {
	case ActionTypeDone,
		ActionTypeDrop,
		"none":
		return true
	default:
		return false
	}
}

func (a *RollDiceAction) NextAction() adventuria.ActionType {
	return ActionTypeRollWheel
}

type RollDiceResult struct {
	Roll        int             `json:"roll"`
	DiceRolls   []int           `json:"dice_roll"`
	CurrentCell adventuria.Cell `json:"current_cell"`
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
	if err != nil {
		return nil, err
	}

	action, err := adventuria.NewActionFromType(a.User(), ActionTypeRollDice)
	if err != nil {
		return nil, err
	}
	action.SetCell(onAfterMoveEvent.CurrentCell.ID())
	action.SetDiceRoll(onAfterMoveEvent.Steps)
	action.SetSeed(rand.Int())
	err = action.Save()
	if err != nil {
		return nil, err
	}

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
		Data: RollDiceResult{
			Roll:        onBeforeRollMoveEvent.N,
			DiceRolls:   diceRolls,
			CurrentCell: onAfterMoveEvent.CurrentCell,
		},
	}, nil
}

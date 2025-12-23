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
	Roll      int                            `json:"roll"`
	DiceRolls []DiceRoll                     `json:"dice_rolls"`
	Path      []*adventuria.OnAfterMoveEvent `json:"path"`
}

type DiceRoll struct {
	Type string `json:"type"`
	Roll int    `json:"roll"`
}

func (a *RollDiceAction) Do(user adventuria.User, _ adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	onBeforeRollEvent := &adventuria.OnBeforeRollEvent{
		Dices: []*adventuria.Dice{adventuria.DiceTypeD6, adventuria.DiceTypeD6},
	}
	res, err := user.OnBeforeRoll().Trigger(onBeforeRollEvent)
	if res != nil && !res.Success {
		return &adventuria.ActionResult{
			Success: false,
			Error:   res.Error,
		}, err
	}
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: failed to trigger onBeforeRollEvent event",
		}, err
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

	res, err = user.OnBeforeRollMove().Trigger(onBeforeRollMoveEvent)
	if res != nil && !res.Success {
		return &adventuria.ActionResult{
			Success: false,
			Error:   res.Error,
		}, err
	}
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: failed to trigger onBeforeRollMoveEvent event",
		}, err
	}

	moveRes, err := user.Move(onBeforeRollMoveEvent.N)
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error",
		}, fmt.Errorf("roll_dice.do(): %w", err)
	}

	user.LastAction().SetType(ActionTypeRollDice)

	onAfterRollEvent := &adventuria.OnAfterRollEvent{
		Dices: onBeforeRollEvent.Dices,
		N:     onBeforeRollMoveEvent.N,
	}
	res, err = user.OnAfterRoll().Trigger(onAfterRollEvent)
	if res != nil && !res.Success {
		return &adventuria.ActionResult{
			Success: false,
			Error:   res.Error,
		}, err
	}
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: failed to trigger onAfterRollEvent event",
		}, err
	}

	return &adventuria.ActionResult{
		Success: true,
		Data: RollDiceResult{
			Roll:      onBeforeRollMoveEvent.N,
			DiceRolls: diceRolls,
			Path:      moveRes,
		},
	}, nil
}

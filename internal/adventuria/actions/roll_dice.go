package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/result"
	"fmt"
)

type RollDiceAction struct {
	adventuria.ActionBase
}

func (a *RollDiceAction) CanDo(ctx adventuria.ActionContext) bool {
	return ctx.User.LastAction().CanMove()
}

type RollDiceResult struct {
	Roll      int                      `json:"roll"`
	DiceRolls []DiceRoll               `json:"dice_rolls"`
	Path      []*adventuria.MoveResult `json:"path"`
}

type DiceRoll struct {
	Type string `json:"type"`
	Roll int    `json:"roll"`
}

func (a *RollDiceAction) Do(ctx adventuria.ActionContext, _ adventuria.ActionRequest) (*result.Result, error) {
	onBeforeRollEvent := &adventuria.OnBeforeRollEvent{
		AppContext: ctx.AppContext,
		Dices:      []*adventuria.Dice{adventuria.DiceTypeD6, adventuria.DiceTypeD6},
	}
	res, err := ctx.User.OnBeforeRoll().Trigger(onBeforeRollEvent)
	if err != nil {
		return res, err
	}
	if res.Failed() {
		return res, err
	}

	onBeforeRollMoveEvent := &adventuria.OnBeforeRollMoveEvent{
		AppContext: ctx.AppContext,
		N:          0,
	}
	diceRolls := make([]DiceRoll, len(onBeforeRollEvent.Dices))
	for i, dice := range onBeforeRollEvent.Dices {
		diceRolls[i] = DiceRoll{
			Type: dice.Type,
			Roll: dice.Roll(),
		}
		onBeforeRollMoveEvent.N += diceRolls[i].Roll
	}

	res, err = ctx.User.OnBeforeRollMove().Trigger(onBeforeRollMoveEvent)
	if err != nil {
		return res, err
	}
	if res.Failed() {
		return res, err
	}

	// we need to save the latest action before ctx.User.Move, because Move() creates a new one
	err = ctx.AppContext.App.Save(ctx.User.LastAction().ProxyRecord())
	if err != nil {
		return result.Err("internal error: can't save action record"),
			fmt.Errorf("roll_dice.do(): %w", err)
	}

	moveRes, err := ctx.User.Move(ctx.AppContext, onBeforeRollMoveEvent.N)
	if err != nil {
		return result.Err("internal error: failed to move to the new cell"),
			fmt.Errorf("roll_dice.do(): %w", err)
	}

	ctx.User.LastAction().SetType(ActionTypeRollDice)

	res, err = ctx.User.OnAfterRoll().Trigger(&adventuria.OnAfterRollEvent{
		AppContext: ctx.AppContext,
		Dices:      onBeforeRollEvent.Dices,
		N:          onBeforeRollMoveEvent.N,
	})
	if err != nil {
		return res, err
	}
	if res.Failed() {
		return res, err
	}

	return result.Ok().WithData(RollDiceResult{
		Roll:      onBeforeRollMoveEvent.N,
		DiceRolls: diceRolls,
		Path:      moveRes,
	}), nil
}

func (a *RollDiceAction) GetVariants(_ adventuria.ActionContext) any {
	return nil
}

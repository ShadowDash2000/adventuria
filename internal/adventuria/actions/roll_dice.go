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
	return ctx.Player.LastAction().CanMove()
}

type RollDiceResult struct {
	Roll      int                      `json:"roll"`
	DiceRolls []DiceRoll               `json:"dice_rolls"`
	Path      []*adventuria.MoveResult `json:"path"`
	From      PositionSnapshot         `json:"from"`
	To        PositionSnapshot         `json:"to"`
	PathSteps []PathStep               `json:"path_steps"`
}

type PositionSnapshot struct {
	WorldId     string `json:"world_id"`
	CellsPassed int    `json:"cells_passed"`
}

type PathStep struct {
	WorldId        string `json:"world_id"`
	WorldSlug      string `json:"world_slug"`
	CellOrder      int    `json:"cell_order"`
	TotalSteps     int    `json:"total_steps"`
	PrevTotalSteps int    `json:"prev_total_steps"`
	Event          string `json:"event"`
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
	res, err := ctx.Player.OnBeforeRoll().Trigger(onBeforeRollEvent)
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

	res, err = ctx.Player.OnBeforeRollMove().Trigger(onBeforeRollMoveEvent)
	if err != nil {
		return res, err
	}
	if res.Failed() {
		return res, err
	}

	// we need to save the latest action before ctx.Player.Move, because Move() creates a new one
	err = ctx.AppContext.App.Save(ctx.Player.LastAction().ProxyRecord())
	if err != nil {
		return result.Err("internal error: can't save action record"),
			fmt.Errorf("roll_dice.do(): %w", err)
	}

	from := PositionSnapshot{
		WorldId:     ctx.Player.Progress().CurrentWorld(),
		CellsPassed: ctx.Player.Progress().CellsPassed(),
	}

	moveRes, err := ctx.Player.Move(ctx.AppContext, onBeforeRollMoveEvent.N)
	if err != nil {
		return result.Err("internal error: failed to move to the new cell"),
			fmt.Errorf("roll_dice.do(): %w", err)
	}

	to := PositionSnapshot{
		WorldId:     ctx.Player.Progress().CurrentWorld(),
		CellsPassed: ctx.Player.Progress().CellsPassed(),
	}

	steps := make([]PathStep, 0, len(moveRes))
	for i, m := range moveRes {
		eventType := "move"
		if i > 0 && moveRes[i-1].CurrentWorld.ID() != m.CurrentWorld.ID() {
			eventType = "world_transition"
		}

		steps = append(steps, PathStep{
			WorldId:        m.CurrentWorld.ID(),
			WorldSlug:      m.CurrentWorld.Slug(),
			CellOrder:      m.CellLocalOrder,
			TotalSteps:     m.TotalSteps,
			PrevTotalSteps: m.PrevTotalSteps,
			Event:          eventType,
		})
	}

	ctx.Player.LastAction().SetType(ActionTypeRollDice)

	res, err = ctx.Player.OnAfterRoll().Trigger(&adventuria.OnAfterRollEvent{
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
		From:      from,
		To:        to,
		PathSteps: steps,
	}), nil
}

func (a *RollDiceAction) GetVariants(_ adventuria.ActionContext) any {
	return nil
}

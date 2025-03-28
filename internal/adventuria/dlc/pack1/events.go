package pack1

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/helper"
)

func WithBaseEvents(g adventuria.Game) adventuria.Game {
	g.Event().On(adventuria.OnAfterReroll, OnAfterRerollStats)
	g.Event().On(adventuria.OnBeforeDrop, OnBeforeDropEffects)
	g.Event().On(adventuria.OnAfterDrop, OnAfterDropStats)
	g.Event().On(adventuria.OnAfterGoToJail, OnAfterGoToJailStats)
	g.Event().On(adventuria.OnBeforeDone, OnBeforeDoneEffects)
	g.Event().On(adventuria.OnAfterDone, OnAfterDoneStats)
	g.Event().On(adventuria.OnBeforeRoll, OnBeforeRollEffects)
	g.Event().On(adventuria.OnBeforeRollMove, OnBeforeRollMoveEffects)
	g.Event().On(adventuria.OnAfterRoll, OnAfterRollStats)
	g.Event().On(adventuria.OnAfterWheelRoll, OnAfterWheelRollStats)
	g.Event().On(adventuria.OnAfterItemUse, OnAfterItemUseStats)
	g.Event().On(adventuria.OnNewLap, OnNewLapItemWheel)

	g.Event().On(adventuria.OnAfterAction, ApplyGenericEffects)

	return g
}

func OnAfterRerollStats(user *adventuria.User, gc *adventuria.GameComponents) error {
	user.Stats.Rerolls++
	return nil
}

func OnBeforeDropEffects(user *adventuria.User, dropEffects *adventuria.DropEffects, gc *adventuria.GameComponents) error {
	effects, _, err := user.Inventory.GetEffects(adventuria.EffectUseOnDrop)
	if err != nil {
		return err
	}

	dropEffects.IsSafeDrop = effects.Effect(EffectTypeSafeDrop).Bool()

	return nil
}

func OnAfterDropStats(user *adventuria.User, gc *adventuria.GameComponents) error {
	user.Stats.Drops++
	return nil
}

func OnAfterGoToJailStats(user *adventuria.User, gc *adventuria.GameComponents) error {
	user.Stats.WasInJail++
	return nil
}

func OnBeforeDoneEffects(user *adventuria.User, doneEffects *adventuria.DoneEffects, gc *adventuria.GameComponents) error {
	effects, _, err := user.Inventory.GetEffects(adventuria.EffectUseOnChooseResult)
	if err != nil {
		return err
	}

	doneEffects.CellPointsDivide = effects.Effect(EffectTypeCellPointsDivide).Int()

	return nil
}

func OnAfterDoneStats(user *adventuria.User, gc *adventuria.GameComponents) error {
	user.Stats.Finished++
	return nil
}

func OnBeforeRollEffects(user *adventuria.User, dicesResult *adventuria.RollDicesResult, gc *adventuria.GameComponents) error {
	effects, _, err := user.Inventory.GetEffects(adventuria.EffectUseOnRoll)
	if err != nil {
		return err
	}

	dicesSrc := adventuria.NewDiceEffectSourceGiver[adventuria.Dice](effects.Effect(EffectTypeChangeDices).Slice())
	dices := dicesSrc.Slice()
	if len(dices) > 0 {
		dicesResult.Dices = dices
	}

	return nil
}

func OnBeforeRollMoveEffects(user *adventuria.User, rollResult *adventuria.RollResult, gc *adventuria.GameComponents) error {
	effects, _, err := user.Inventory.GetEffects(adventuria.EffectUseOnRoll)
	if err != nil {
		return err
	}

	diceMultiplier := effects.Effect(EffectTypeDiceMultiplier).Int()
	if diceMultiplier > 0 {
		rollResult.N *= diceMultiplier
	}

	diceIncrement := effects.Effect(EffectTypeDiceIncrement).Int()
	rollResult.N += diceIncrement

	rollReverse := effects.Effect(EffectTypeRollReverse).Bool()
	if rollReverse {
		rollResult.N *= -1
	}

	return nil
}

func OnAfterRollStats(user *adventuria.User, rollResult *adventuria.RollResult, gc *adventuria.GameComponents) error {
	user.Stats.DiceRolls++
	if rollResult.N > user.Stats.MaxDiceRoll {
		user.Stats.MaxDiceRoll = rollResult.N
	}
	return nil
}

func OnAfterWheelRollStats(user *adventuria.User, gc *adventuria.GameComponents) error {
	user.Stats.WheelRolled++
	return nil
}

func OnAfterItemUseStats(user *adventuria.User, gc *adventuria.GameComponents) error {
	user.Stats.ItemsUsed++
	return nil
}

func ApplyGenericEffects(user *adventuria.User, event string, gc *adventuria.GameComponents) error {
	effects, _, err := user.Inventory.GetEffects(event)
	if err != nil {
		return err
	}

	pointsIncrement := effects.Effect(EffectTypePointsIncrement).Int()
	if pointsIncrement != 0 {
		user.SetPoints(user.Points() + pointsIncrement)
	}

	timerIncrement := effects.Effect(EffectTypeTimerIncrement).Int()
	if timerIncrement != 0 {
		err := user.Timer.AddSecondsTimeLimit(timerIncrement)
		if err != nil {
			return err
		}
	}

	jailEscape := effects.Effect(EffectTypeJailEscape).Bool()
	if jailEscape {
		user.SetIsInJail(false)
		user.SetDropsInARow(0)
	}

	dropInventory := effects.Effect(EffectTypeDropInventory).Bool()
	if dropInventory {
		err := user.Inventory.DropInventory()
		if err != nil {
			return err
		}
	}

	cellTypes := effects.Effect(EffectTypeTeleportToRandomCellByTypes).Slice()
	if len(cellTypes) > 0 {
		cells := gc.Cells.GetAllByTypes(cellTypes)
		if currentCell, ok := user.CurrentCell(); ok {
			cells = adventuria.FilterByField(cells, []string{currentCell.Id}, func(cell *adventuria.Cell) string {
				return cell.Id
			})
		}
		cell := helper.RandomItemFromSlice(cells)
		err := user.MoveToCellId(cell.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func OnNewLapItemWheel(user *adventuria.User, laps int, gc *adventuria.GameComponents) error {
	// Every lap gives one item wheel
	user.SetItemWheelsCount(user.ItemWheelsCount() + laps)

	return nil
}

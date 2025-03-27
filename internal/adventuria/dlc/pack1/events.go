package pack1

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/helper"
)

func WithBaseEvents(g adventuria.Game) adventuria.Game {
	g.Event().On(adventuria.OnAfterChooseGame, OnAfterChooseGameEffects)
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
	g.Event().On(adventuria.OnAfterItemRoll, OnAfterItemRollEffects)
	g.Event().On(adventuria.OnAfterItemUse, OnAfterItemUseEffects)
	g.Event().On(adventuria.OnAfterItemUse, OnAfterItemUseStats)
	g.Event().On(adventuria.OnNewLap, OnNewLapItemWheel)

	return g
}

func OnAfterChooseGameEffects(user *adventuria.User, gc *adventuria.GameComponents) error {
	effects, _, err := user.Inventory.GetEffects(adventuria.EffectUseOnChooseGame)
	if err != nil {
		return err
	}

	err = ApplyGenericEffects(effects, user, gc)
	if err != nil {
		return err
	}

	return nil
}

func OnAfterRerollStats(user *adventuria.User) error {
	user.Stats.Rerolls++
	return nil
}

func OnBeforeDropEffects(user *adventuria.User, dropEffects *adventuria.DropEffects, gc *adventuria.GameComponents) error {
	effects, _, err := user.Inventory.GetEffects(adventuria.EffectUseOnDrop)
	if err != nil {
		return err
	}

	dropEffects.IsSafeDrop = effects.Effect(EffectTypeSafeDrop).Bool()

	err = ApplyGenericEffects(effects, user, gc)
	if err != nil {
		return err
	}

	return nil
}

func OnAfterDropStats(user *adventuria.User) error {
	user.Stats.Drops++
	return nil
}

func OnAfterGoToJailStats(user *adventuria.User) error {
	user.Stats.WasInJail++
	return nil
}

func OnBeforeDoneEffects(user *adventuria.User, doneEffects *adventuria.DoneEffects, gc *adventuria.GameComponents) error {
	effects, _, err := user.Inventory.GetEffects(adventuria.EffectUseOnChooseResult)
	if err != nil {
		return err
	}

	doneEffects.CellPointsDivide = effects.Effect(EffectTypeCellPointsDivide).Int()

	err = ApplyGenericEffects(effects, user, gc)
	if err != nil {
		return err
	}

	return nil
}

func OnAfterDoneStats(user *adventuria.User) error {
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

	err = ApplyGenericEffects(effects, user, gc)
	if err != nil {
		return err
	}

	return nil
}

func OnBeforeRollMoveEffects(user *adventuria.User, rollResult *adventuria.RollResult) error {
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

func OnAfterRollStats(user *adventuria.User, rollResult *adventuria.RollResult) error {
	user.Stats.DiceRolls++
	if rollResult.N > user.Stats.MaxDiceRoll {
		user.Stats.MaxDiceRoll = rollResult.N
	}
	return nil
}

func OnAfterWheelRollStats(user *adventuria.User) error {
	user.Stats.WheelRolled++
	return nil
}

func OnAfterItemRollEffects(user *adventuria.User, gc *adventuria.GameComponents) error {
	effects, _, err := user.Inventory.GetEffects(adventuria.EffectUseOnRollItem)
	if err != nil {
		return err
	}

	err = ApplyGenericEffects(effects, user, gc)
	if err != nil {
		return err
	}

	return nil
}

func OnAfterItemUseEffects(user *adventuria.User, gc *adventuria.GameComponents) error {
	effects, _, err := user.Inventory.GetEffects(adventuria.EffectUseInstant)
	if err != nil {
		return err
	}

	err = ApplyGenericEffects(effects, user, gc)
	if err != nil {
		return err
	}

	return nil
}

func OnAfterItemUseStats(user *adventuria.User) error {
	user.Stats.ItemsUsed++
	return nil
}

func ApplyGenericEffects(effects *adventuria.Effects, user *adventuria.User, gc *adventuria.GameComponents) error {
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

func OnNewLapItemWheel(user *adventuria.User, laps int) error {
	// Every lap gives one item wheel
	user.SetItemWheelsCount(user.ItemWheelsCount() + laps)

	currentCell, _ := user.CurrentCell()

	// If we moved to item cell type, we need to add one additional item roll
	if currentCell.Type() == adventuria.CellTypeItem {
		user.SetItemWheelsCount(user.ItemWheelsCount() + 1)
	}

	return nil
}

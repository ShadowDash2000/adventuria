package pack1

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/helper"
)

func WithBaseEvents(g adventuria.Game) adventuria.Game {
	adventuria.GameEvent.On(adventuria.OnAfterReroll, OnAfterRerollStats)
	adventuria.GameEvent.On(adventuria.OnBeforeDrop, OnBeforeDropEffects)
	adventuria.GameEvent.On(adventuria.OnAfterDrop, OnAfterDropStats)
	adventuria.GameEvent.On(adventuria.OnAfterGoToJail, OnAfterGoToJailStats)
	adventuria.GameEvent.On(adventuria.OnBeforeDone, OnBeforeDoneEffects)
	adventuria.GameEvent.On(adventuria.OnAfterDone, OnAfterDoneStats)
	adventuria.GameEvent.On(adventuria.OnBeforeRoll, OnBeforeRollEffects)
	adventuria.GameEvent.On(adventuria.OnBeforeRollMove, OnBeforeRollMoveEffects)
	adventuria.GameEvent.On(adventuria.OnAfterRoll, OnAfterRollStats)
	adventuria.GameEvent.On(adventuria.OnAfterWheelRoll, OnAfterWheelRollStats)
	adventuria.GameEvent.On(adventuria.OnAfterItemUse, OnAfterItemUseStats)

	// item wheels
	adventuria.GameEvent.On(adventuria.OnNewLap, OnNewLapItemWheel)
	adventuria.GameEvent.On(adventuria.OnAfterDone, GiveOneItemWheel)

	adventuria.GameEvent.On(adventuria.OnAfterAction, ApplyGenericEffects)

	return g
}

func OnAfterRerollStats(e adventuria.EventFields) error {
	e.User().Stats.Rerolls++
	return nil
}

func OnBeforeDropEffects(e adventuria.EventFields) error {
	effects := e.Effects(adventuria.EffectUseOnDrop)
	fields := e.Fields().(*adventuria.OnBeforeDropFields)

	fields.IsSafeDrop = effects.Effect(EffectTypeSafeDrop).Bool()

	return nil
}

func OnAfterDropStats(e adventuria.EventFields) error {
	e.User().Stats.Drops++
	return nil
}

func OnAfterGoToJailStats(e adventuria.EventFields) error {
	e.User().Stats.WasInJail++
	return nil
}

func OnBeforeDoneEffects(e adventuria.EventFields) error {
	effects := e.Effects(adventuria.EffectUseOnDone)
	fields := e.Fields().(*adventuria.OnBeforeDoneFields)

	fields.CellPointsDivide = effects.Effect(EffectTypeCellPointsDivide).Int()

	return nil
}

func OnAfterDoneStats(e adventuria.EventFields) error {
	e.User().Stats.Finished++
	return nil
}

func OnBeforeRollEffects(e adventuria.EventFields) error {
	effects := e.Effects(adventuria.EffectUseOnRoll)
	fields := e.Fields().(*adventuria.OnBeforeRollFields)

	dicesSrc := adventuria.NewDiceEffectSourceGiver(effects.Effect(EffectTypeChangeDices).Slice())
	dices := dicesSrc.Slice()
	if len(dices) > 0 {
		fields.Dices = dices
	}

	return nil
}

func OnBeforeRollMoveEffects(e adventuria.EventFields) error {
	effects := e.Effects(adventuria.EffectUseOnRoll)
	fields := e.Fields().(*adventuria.OnBeforeRollMoveFields)

	diceMultiplier := effects.Effect(EffectTypeDiceMultiplier).Int()
	if diceMultiplier > 0 {
		fields.N *= diceMultiplier
	}

	diceIncrement := effects.Effect(EffectTypeDiceIncrement).Int()
	fields.N += diceIncrement

	rollReverse := effects.Effect(EffectTypeRollReverse).Bool()
	if rollReverse {
		fields.N *= -1
	}

	return nil
}

func OnAfterRollStats(e adventuria.EventFields) error {
	fields := e.Fields().(*adventuria.OnAfterRollFields)

	e.User().Stats.DiceRolls++

	if fields.N > e.User().Stats.MaxDiceRoll {
		e.User().Stats.MaxDiceRoll = fields.N
	}
	return nil
}

func OnAfterWheelRollStats(e adventuria.EventFields) error {
	e.User().Stats.WheelRolled++
	return nil
}

func OnAfterItemUseStats(e adventuria.EventFields) error {
	e.User().Stats.ItemsUsed++
	return nil
}

func ApplyGenericEffects(e adventuria.EventFields) error {
	fields := e.Fields().(*adventuria.OnAfterActionFields)
	effects := e.Effects(fields.Event)

	pointsIncrement := effects.Effect(EffectTypePointsIncrement).Int()
	if pointsIncrement != 0 {
		e.User().SetPoints(e.User().Points() + pointsIncrement)
	}

	timerIncrement := effects.Effect(EffectTypeTimerIncrement).Int()
	if timerIncrement != 0 {
		err := e.User().Timer.AddSecondsTimeLimit(timerIncrement)
		if err != nil {
			return err
		}
	}

	jailEscape := effects.Effect(EffectTypeJailEscape).Bool()
	if jailEscape {
		e.User().SetIsInJail(false)
		e.User().SetDropsInARow(0)
	}

	dropInventory := effects.Effect(EffectTypeDropInventory).Bool()
	if dropInventory {
		err := e.User().Inventory.DropInventory()
		if err != nil {
			return err
		}
	}

	cellTypesSrc := adventuria.NewCellTypeSourceGiver(effects.Effect(EffectTypeTeleportToRandomCellByTypes).Slice())
	cellTypes := cellTypesSrc.Slice()
	if len(cellTypes) > 0 {
		cells := adventuria.GameCells.GetAllByTypes(cellTypes)
		if currentCell, ok := e.User().CurrentCell(); ok {
			cells = helper.FilterByField(cells, []string{currentCell.ID()}, func(cell adventuria.Cell) string {
				return cell.ID()
			})
		}
		cell := helper.RandomItemFromSlice(cells)
		err := e.User().MoveToCellId(cell.ID())
		if err != nil {
			return err
		}
	}

	return nil
}

func GiveOneItemWheel(e adventuria.EventFields) error {
	e.User().SetItemWheelsCount(e.User().ItemWheelsCount() + 1)

	return nil
}

func OnNewLapItemWheel(e adventuria.EventFields) error {
	// Every lap gives one item wheel
	fields := e.Fields().(*adventuria.OnNewLapFields)

	e.User().SetItemWheelsCount(e.User().ItemWheelsCount() + fields.Laps)

	return nil
}

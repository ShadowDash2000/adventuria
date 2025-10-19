package effects

import "adventuria/internal/adventuria"

func WithBaseEffects() {
	adventuria.RegisterEffects(map[string]adventuria.EffectCreator{
		"pointsIncrement":             adventuria.NewEffect(&PointsIncrementEffect{}),
		"jailEscape":                  adventuria.NewEffect(&JailEscapeEffect{}),
		"diceMultiplier":              adventuria.NewEffect(&DiceMultiplierEffect{}),
		"diceIncrement":               adventuria.NewEffect(&DiceIncrementEffect{}),
		"changeDices":                 nil,
		"safeDrop":                    adventuria.NewEffect(&SafeDropEffect{}),
		"timerIncrement":              adventuria.NewEffect(&TimerIncrementEffect{}),
		"rollReverse":                 adventuria.NewEffect(&RollReverseEffect{}),
		"dropInventory":               adventuria.NewEffect(&DropInventoryEffect{}),
		"cellPointsDivide":            adventuria.NewEffect(&CellPointsDivideEffect{}),
		"teleportToRandomCellByTypes": nil,
		"teleportToRandomCellByIds":   adventuria.NewEffect(&TeleportToRandomCellByIdEffect{}),
	})

	adventuria.RegisterPersistentEffects(map[string]adventuria.PersistentEffectCreator{
		"stats":                 adventuria.NewPersistentEffect(&StatsEffect{}),
		"give_wheel_on_done":    adventuria.NewPersistentEffect(&GiveWheelOnDoneEffect{}),
		"give_wheel_on_new_lap": adventuria.NewPersistentEffect(&GiveWheelOnNewLapEffect{}),
	})
}

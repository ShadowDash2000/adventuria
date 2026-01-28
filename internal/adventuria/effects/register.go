package effects

import "adventuria/internal/adventuria"

func WithBaseEffects() {
	adventuria.RegisterEffects(map[string]adventuria.EffectCreator{
		"nothing":                     adventuria.NewEffect(&NothingEffect{}),
		"pointsIncrement":             adventuria.NewEffect(&PointsIncrementEffect{}),
		"jailEscape":                  adventuria.NewEffect(&JailEscapeEffect{}),
		"diceMultiplier":              adventuria.NewEffect(&DiceMultiplierEffect{}),
		"diceIncrement":               adventuria.NewEffect(&DiceIncrementEffect{}),
		"changeDices":                 adventuria.NewEffect(&ChangeDicesEffect{}),
		"safeDrop":                    adventuria.NewEffect(&SafeDropEffect{}),
		"timerIncrement":              adventuria.NewEffect(&TimerIncrementEffect{}),
		"rollReverse":                 adventuria.NewEffect(&RollReverseEffect{}),
		"dropInventory":               adventuria.NewEffect(&DropInventoryEffect{}),
		"cellPointsDivide":            adventuria.NewEffect(&CellPointsDivideEffect{}),
		"teleportToRandomCellByName":  adventuria.NewEffect(&TeleportToRandomCellByNameEffect{}),
		"teleportToRandomCellByType":  adventuria.NewEffect(&TeleportToRandomCellByTypeEffect{}),
		"teleportToClosestCellByName": adventuria.NewEffect(&TeleportToClosestCellByNameEffect{}),
		"changeMinGamePrice":          adventuria.NewEffect(&ChangeMinGamePriceEffect{}),
		"changeMaxGamePrice":          adventuria.NewEffect(&ChangeMaxGamePriceEffect{}),
		"noTimeLimit":                 adventuria.NewEffect(&NoTimeLimitEffect{}),
		"addGameGenre":                adventuria.NewEffect(&AddGameGenreEffect{}),
		"replaceDiceRoll":             adventuria.NewEffect(&ReplaceDiceRollEffect{}),
		"addRandomItemToInventory":    adventuria.NewEffect(&AddRandomItemToInventoryEffect{}),
		"goToJail":                    adventuria.NewEffect(&GoToJailEffect{}),
		"changeGameById":              adventuria.NewEffect(&ChangeGameByIdEffect{}),
		"chooseGame":                  adventuria.NewEffect(&ChooseGameEffect{}),
		"dropPointsDivide":            adventuria.NewEffect(&DropPointsDivideEffect{}),
		"returnToPrevCell":            adventuria.NewEffect(&ReturnToPrevCellEffect{}),
		"noCoinsForDone":              adventuria.NewEffect(&NoCoinsForDoneEffect{}),
		"dropBlock":                   adventuria.NewEffect(&DropBlockedEffect{}),
		"stayOnCellAfterDone":         adventuria.NewEffect(&StayOnCellAfterDoneEffect{}),
		"debuffBlock":                 adventuria.NewEffect(&DebuffBlockEffect{}),
		"paidMovementInRadius":        adventuria.NewEffect(&PaidMovementInRadiusEffect{}),
	})

	adventuria.RegisterPersistentEffects(map[string]adventuria.PersistentEffectCreator{
		"stats":                 adventuria.NewPersistentEffect(&StatsEffect{}),
		"give_wheel_on_done":    adventuria.NewPersistentEffect(&GiveWheelOnDoneEffect{}),
		"give_wheel_on_new_lap": adventuria.NewPersistentEffect(&GiveWheelOnNewLapEffect{}),
	})
}

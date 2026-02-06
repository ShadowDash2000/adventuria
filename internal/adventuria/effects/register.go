package effects

import "adventuria/internal/adventuria"

func WithBaseEffects() {
	adventuria.RegisterEffects(map[string]adventuria.EffectCreator{
		"nothing":                     adventuria.NewEffect(func() adventuria.Effect { return &NothingEffect{} }),
		"pointsIncrement":             adventuria.NewEffect(func() adventuria.Effect { return &PointsIncrementEffect{} }),
		"jailEscape":                  adventuria.NewEffect(func() adventuria.Effect { return &JailEscapeEffect{} }),
		"diceMultiplier":              adventuria.NewEffect(func() adventuria.Effect { return &DiceMultiplierEffect{} }),
		"diceIncrement":               adventuria.NewEffect(func() adventuria.Effect { return &DiceIncrementEffect{} }),
		"changeDices":                 adventuria.NewEffect(func() adventuria.Effect { return &ChangeDicesEffect{} }),
		"safeDrop":                    adventuria.NewEffect(func() adventuria.Effect { return &SafeDropEffect{} }),
		"timerIncrement":              adventuria.NewEffect(func() adventuria.Effect { return &TimerIncrementEffect{} }),
		"rollReverse":                 adventuria.NewEffect(func() adventuria.Effect { return &RollReverseEffect{} }),
		"dropInventory":               adventuria.NewEffect(func() adventuria.Effect { return &DropInventoryEffect{} }),
		"cellPointsDivide":            adventuria.NewEffect(func() adventuria.Effect { return &CellPointsDivideEffect{} }),
		"teleportToRandomCellByName":  adventuria.NewEffect(func() adventuria.Effect { return &TeleportToRandomCellByNameEffect{} }),
		"teleportToRandomCellByType":  adventuria.NewEffect(func() adventuria.Effect { return &TeleportToRandomCellByTypeEffect{} }),
		"teleportToClosestCellByName": adventuria.NewEffect(func() adventuria.Effect { return &TeleportToClosestCellByNameEffect{} }),
		"changeMinGamePrice":          adventuria.NewEffect(func() adventuria.Effect { return &ChangeMinGamePriceEffect{} }),
		"changeMaxGamePrice":          adventuria.NewEffect(func() adventuria.Effect { return &ChangeMaxGamePriceEffect{} }),
		"noTimeLimit":                 adventuria.NewEffect(func() adventuria.Effect { return &NoTimeLimitEffect{} }),
		"addGameGenre":                adventuria.NewEffect(func() adventuria.Effect { return &AddGameGenreEffect{} }),
		"replaceDiceRoll":             adventuria.NewEffect(func() adventuria.Effect { return &ReplaceDiceRollEffect{} }),
		"addRandomItemToInventory":    adventuria.NewEffect(func() adventuria.Effect { return &AddRandomItemToInventoryEffect{} }),
		"goToJail":                    adventuria.NewEffect(func() adventuria.Effect { return &GoToJailEffect{} }),
		"changeGameById":              adventuria.NewEffect(func() adventuria.Effect { return &ChangeGameByIdEffect{} }),
		"chooseGame":                  adventuria.NewEffect(func() adventuria.Effect { return &ChooseGameEffect{} }),
		"dropPointsDivide":            adventuria.NewEffect(func() adventuria.Effect { return &DropPointsDivideEffect{} }),
		"returnToPrevCell":            adventuria.NewEffect(func() adventuria.Effect { return &ReturnToPrevCellEffect{} }),
		"noCoinsForDone":              adventuria.NewEffect(func() adventuria.Effect { return &NoCoinsForDoneEffect{} }),
		"dropBlock":                   adventuria.NewEffect(func() adventuria.Effect { return &DropBlockedEffect{} }),
		"rerollBlock":                 adventuria.NewEffect(func() adventuria.Effect { return &RerollBlockedEffect{} }),
		"stayOnCellAfterDone":         adventuria.NewEffect(func() adventuria.Effect { return &StayOnCellAfterDoneEffect{} }),
		"debuffBlock":                 adventuria.NewEffect(func() adventuria.Effect { return &DebuffBlockEffect{} }),
		"paidMovementInRadius":        adventuria.NewEffect(func() adventuria.Effect { return &PaidMovementInRadiusEffect{} }),
		"coinsIncrement":              adventuria.NewEffect(func() adventuria.Effect { return &CoinsIncrementEffect{} }),
		"discountPriceDivide":         adventuria.NewEffect(func() adventuria.Effect { return &DiscountPriceDivideEffect{} }),
		"coinsForAll":                 adventuria.NewEffect(func() adventuria.Effect { return &CoinsForAllEffect{} }),
	})

	adventuria.RegisterPersistentEffects(map[string]adventuria.PersistentEffectCreator{
		"stats":                 adventuria.NewPersistentEffect(func() adventuria.PersistentEffect { return &StatsEffect{} }),
		"give_wheel_on_done":    adventuria.NewPersistentEffect(func() adventuria.PersistentEffect { return &GiveWheelOnDoneEffect{} }),
		"give_wheel_on_new_lap": adventuria.NewPersistentEffect(func() adventuria.PersistentEffect { return &GiveWheelOnNewLapEffect{} }),
	})
}

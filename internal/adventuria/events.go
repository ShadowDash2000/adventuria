package adventuria

const (
	OnAfterChooseGame = "OnAfterChooseGame"
	OnAfterReroll     = "OnAfterReroll"
	OnBeforeDrop      = "OnBeforeDrop"
	OnAfterDrop       = "OnAfterDrop"
	OnAfterGoToJail   = "OnAfterGoToJail"
	OnBeforeDone      = "OnBeforeDone"
	OnAfterDone       = "OnAfterDone"
	OnBeforeRoll      = "OnBeforeRoll"
	OnAfterRoll       = "OnAfterRoll"
	OnAfterWheelRoll  = "OnAfterWheelRoll"
	OnAfterItemRoll   = "OnAfterItemRoll"
	OnAfterItemUse    = "OnAfterItemUse"
)

func (g *Game) bindEvents() {
	g.GC.event.On(OnAfterChooseGame, OnAfterChooseGameEffects)
	g.GC.event.On(OnAfterReroll, OnAfterRerollStats)
	g.GC.event.On(OnBeforeDrop, OnBeforeDropEffects)
	g.GC.event.On(OnAfterDrop, OnAfterDropStats)
	g.GC.event.On(OnAfterGoToJail, OnAfterGoToJailStats)
	g.GC.event.On(OnBeforeDone, OnBeforeDoneEffects)
	g.GC.event.On(OnAfterDone, OnAfterDoneStats)
	g.GC.event.On(OnBeforeRoll, OnBeforeRollEffects)
	g.GC.event.On(OnAfterRoll, OnAfterRollStats)
	g.GC.event.On(OnAfterWheelRoll, OnAfterWheelRollStats)
	g.GC.event.On(OnAfterItemRoll, OnAfterItemRollEffects)
	g.GC.event.On(OnAfterItemUse, OnAfterItemUseEffects)
	g.GC.event.On(OnAfterItemUse, OnAfterItemUseStats)
}

func OnAfterChooseGameEffects(user *User) error {
	effects, err := user.Inventory.ApplyEffects(ItemUseOnChooseGame)
	if err != nil {
		return err
	}

	err = ApplyGenericEffects(effects, user)
	if err != nil {
		return err
	}

	return nil
}

func OnAfterRerollStats(user *User) error {
	user.Stats.Rerolls++
	return user.Save()
}

func OnBeforeDropEffects(user *User, dropEffects *DropEffects) error {
	effects, err := user.Inventory.ApplyEffects(ItemUseOnDrop)
	if err != nil {
		return err
	}

	dropEffects.IsSafeDrop = effects.IsSafeDrop

	err = ApplyGenericEffects(effects, user)
	if err != nil {
		return err
	}

	return nil
}

func OnAfterDropStats(user *User) error {
	user.Stats.Drops++
	return user.Save()
}

func OnAfterGoToJailStats(user *User) error {
	user.Stats.WasInJail++
	return user.Save()
}

func OnBeforeDoneEffects(user *User, doneEffects *DoneEffects) error {
	effects, err := user.Inventory.ApplyEffects(ItemUseOnChooseResult)
	if err != nil {
		return err
	}

	doneEffects.CellPointsDivide = effects.CellPointsDivide

	err = ApplyGenericEffects(effects, user)
	if err != nil {
		return err
	}

	return nil
}

func OnAfterDoneStats(user *User) error {
	user.Stats.Finished++
	return user.Save()
}

func OnBeforeRollEffects(user *User, rollEffects *RollEffects) error {
	effects, err := user.Inventory.ApplyEffects(ItemUseOnRoll)
	if err != nil {
		return err
	}

	rollEffects.DiceMultiplier = effects.DiceMultiplier
	rollEffects.DiceIncrement = effects.DiceIncrement
	rollEffects.Dices = effects.Dices
	rollEffects.RollReverse = effects.RollReverse

	err = ApplyGenericEffects(effects, user)
	if err != nil {
		return err
	}

	return nil
}

func OnAfterRollStats(user *User, roll *RollResult) error {
	user.Stats.DiceRolls++
	if roll.n > user.Stats.MaxDiceRoll {
		user.Stats.MaxDiceRoll = roll.n
	}
	return user.Save()
}

func OnAfterWheelRollStats(user *User) error {
	user.Stats.WheelRolled++
	return user.Save()
}

func OnAfterItemRollEffects(user *User) error {
	effects, err := user.Inventory.ApplyEffects(ItemUseOnRollItem)
	if err != nil {
		return err
	}

	err = ApplyGenericEffects(effects, user)
	if err != nil {
		return err
	}

	return nil
}

func OnAfterItemUseEffects(user *User) error {
	effects, err := user.Inventory.ApplyEffects(ItemUseInstant)
	if err != nil {
		return err
	}

	err = ApplyGenericEffects(effects, user)
	if err != nil {
		return err
	}

	return nil
}

func OnAfterItemUseStats(user *User) error {
	user.Stats.ItemsUsed++
	return user.Save()
}

func ApplyGenericEffects(effects *Effects, user *User) error {
	if effects.PointsIncrement != 0 {
		user.Set("points", user.Points()+effects.PointsIncrement)
		err := user.Save()
		if err != nil {
			return err
		}
	}

	if effects.TimerIncrement != 0 {
		err := user.Timer.AddSecondsTimeLimit(effects.TimerIncrement)
		if err != nil {
			return err
		}
	}

	if effects.JailEscape {
		user.Set("isInJail", false)
		err := user.Save()
		if err != nil {
			return err
		}
	}

	if effects.DropInventory {
		err := user.Inventory.DropInventory()
		if err != nil {
			return err
		}
	}

	return nil
}

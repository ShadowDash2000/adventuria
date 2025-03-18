package adventuria

func WithBaseEvents(g Game) Game {
	g.Event().On(OnAfterChooseGame, OnAfterChooseGameEffects)
	g.Event().On(OnAfterReroll, OnAfterRerollStats)
	g.Event().On(OnBeforeDrop, OnBeforeDropEffects)
	g.Event().On(OnAfterDrop, OnAfterDropStats)
	g.Event().On(OnAfterGoToJail, OnAfterGoToJailStats)
	g.Event().On(OnBeforeDone, OnBeforeDoneEffects)
	g.Event().On(OnAfterDone, OnAfterDoneStats)
	g.Event().On(OnBeforeRoll, OnBeforeRollEffects)
	g.Event().On(OnBeforeRollMove, OnBeforeRollMoveEffects)
	g.Event().On(OnAfterRoll, OnAfterRollStats)
	g.Event().On(OnAfterWheelRoll, OnAfterWheelRollStats)
	g.Event().On(OnAfterItemRoll, OnAfterItemRollEffects)
	g.Event().On(OnAfterItemUse, OnAfterItemUseEffects)
	g.Event().On(OnAfterItemUse, OnAfterItemUseStats)

	return g
}

func OnAfterChooseGameEffects(user *User) error {
	effects, _, err := user.Inventory.GetEffects(EffectUseOnChooseGame)
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
	return nil
}

func OnBeforeDropEffects(user *User, dropEffects *Effects) error {
	effects, _, err := user.Inventory.GetEffects(EffectUseOnDrop)
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
	return nil
}

func OnAfterGoToJailStats(user *User) error {
	user.Stats.WasInJail++
	return nil
}

func OnBeforeDoneEffects(user *User, doneEffects *Effects) error {
	effects, _, err := user.Inventory.GetEffects(EffectUseOnChooseResult)
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
	return nil
}

func OnBeforeRollEffects(user *User, dicesResult *RollDicesResult) error {
	effects, _, err := user.Inventory.GetEffects(EffectUseOnRoll)
	if err != nil {
		return err
	}

	if len(effects.Dices) > 0 {
		dicesResult.Dices = effects.Dices
	}

	err = ApplyGenericEffects(effects, user)
	if err != nil {
		return err
	}

	return nil
}

func OnBeforeRollMoveEffects(user *User, rollResult *RollResult) error {
	effects, _, err := user.Inventory.GetEffects(EffectUseOnRoll)
	if err != nil {
		return err
	}

	if effects.DiceMultiplier > 0 {
		rollResult.N *= effects.DiceMultiplier
	}
	rollResult.N += effects.DiceIncrement

	if effects.RollReverse {
		rollResult.N *= -1
	}

	return nil
}

func OnAfterRollStats(user *User, rollResult *RollResult) error {
	user.Stats.DiceRolls++
	if rollResult.N > user.Stats.MaxDiceRoll {
		user.Stats.MaxDiceRoll = rollResult.N
	}
	return nil
}

func OnAfterWheelRollStats(user *User) error {
	user.Stats.WheelRolled++
	return nil
}

func OnAfterItemRollEffects(user *User) error {
	effects, _, err := user.Inventory.GetEffects(EffectUseOnRollItem)
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
	effects, _, err := user.Inventory.GetEffects(EffectUseInstant)
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
	return nil
}

func ApplyGenericEffects(effects *Effects, user *User) error {
	if effects.PointsIncrement != 0 {
		user.SetPoints(user.Points() + effects.PointsIncrement)
	}

	if effects.TimerIncrement != 0 {
		err := user.Timer.AddSecondsTimeLimit(effects.TimerIncrement)
		if err != nil {
			return err
		}
	}

	if effects.JailEscape {
		user.SetIsInJail(false)
	}

	if effects.DropInventory {
		err := user.Inventory.DropInventory()
		if err != nil {
			return err
		}
	}

	return nil
}

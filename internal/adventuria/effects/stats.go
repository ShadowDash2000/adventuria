package effects

import (
	"adventuria/internal/adventuria"
)

type StatsEffect struct {
	adventuria.PersistentEffectBase
}

func (ef *StatsEffect) Subscribe() {
	ef.PoolUnsubscribers(
		ef.User().OnAfterReroll().BindFunc(func(e *adventuria.OnAfterRerollEvent) error {
			ef.User().Stats().Rerolls++
			return e.Next()
		}),
		ef.User().OnAfterDrop().BindFunc(func(e *adventuria.OnAfterDropEvent) error {
			ef.User().Stats().Drops++
			return e.Next()
		}),
		ef.User().OnAfterGoToJail().BindFunc(func(e *adventuria.OnAfterGoToJailEvent) error {
			ef.User().Stats().WasInJail++
			return e.Next()
		}),
		ef.User().OnAfterDone().BindFunc(func(e *adventuria.OnAfterDoneEvent) error {
			ef.User().Stats().Finished++
			return e.Next()
		}),
		ef.User().OnAfterRoll().BindFunc(func(e *adventuria.OnAfterRollEvent) error {
			ef.User().Stats().DiceRolls++
			if e.N > ef.User().Stats().MaxDiceRoll {
				ef.User().Stats().MaxDiceRoll = e.N
			}
			return e.Next()
		}),
		ef.User().OnAfterWheelRoll().BindFunc(func(e *adventuria.OnAfterWheelRollEvent) error {
			ef.User().Stats().WheelRolled++
			return e.Next()
		}),
		ef.User().OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) error {
			ef.User().Stats().ItemsUsed++
			return e.Next()
		}),
	)
}

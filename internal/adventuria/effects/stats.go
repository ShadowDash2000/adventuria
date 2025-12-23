package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type StatsEffect struct {
	adventuria.PersistentEffectBase
}

func (ef *StatsEffect) Subscribe(user adventuria.User) []event.Unsubscribe {
	return []event.Unsubscribe{
		user.OnAfterReroll().BindFunc(func(e *adventuria.OnAfterRerollEvent) (*event.Result, error) {
			user.Stats().Rerolls++
			return e.Next()
		}),
		user.OnAfterDrop().BindFunc(func(e *adventuria.OnAfterDropEvent) (*event.Result, error) {
			user.Stats().Drops++
			return e.Next()
		}),
		user.OnAfterGoToJail().BindFunc(func(e *adventuria.OnAfterGoToJailEvent) (*event.Result, error) {
			user.Stats().WasInJail++
			return e.Next()
		}),
		user.OnAfterDone().BindFunc(func(e *adventuria.OnAfterDoneEvent) (*event.Result, error) {
			user.Stats().Finished++
			return e.Next()
		}),
		user.OnAfterRoll().BindFunc(func(e *adventuria.OnAfterRollEvent) (*event.Result, error) {
			user.Stats().DiceRolls++
			if e.N > user.Stats().MaxDiceRoll {
				user.Stats().MaxDiceRoll = e.N
			}
			return e.Next()
		}),
		user.OnAfterWheelRoll().BindFunc(func(e *adventuria.OnAfterWheelRollEvent) (*event.Result, error) {
			user.Stats().WheelRolled++
			return e.Next()
		}),
		user.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*event.Result, error) {
			user.Stats().ItemsUsed++
			return e.Next()
		}),
	}
}

package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
)

type StatsEffect struct{}

func (ef *StatsEffect) Subscribe(user adventuria.User) []event.Unsubscribe {
	return []event.Unsubscribe{
		user.OnAfterReroll().BindFunc(func(e *adventuria.OnAfterRerollEvent) (*result.Result, error) {
			user.Stats().Rerolls++
			return e.Next()
		}),
		user.OnAfterDrop().BindFunc(func(e *adventuria.OnAfterDropEvent) (*result.Result, error) {
			user.Stats().Drops++
			return e.Next()
		}),
		user.OnAfterGoToJail().BindFunc(func(e *adventuria.OnAfterGoToJailEvent) (*result.Result, error) {
			user.Stats().WasInJail++
			return e.Next()
		}),
		user.OnAfterDone().BindFunc(func(e *adventuria.OnAfterDoneEvent) (*result.Result, error) {
			user.Stats().Finished++
			return e.Next()
		}),
		user.OnAfterRoll().BindFunc(func(e *adventuria.OnAfterRollEvent) (*result.Result, error) {
			user.Stats().DiceRolls++
			if e.N > user.Stats().MaxDiceRoll {
				user.Stats().MaxDiceRoll = e.N
			}
			return e.Next()
		}),
		user.OnAfterWheelRoll().BindFunc(func(e *adventuria.OnAfterWheelRollEvent) (*result.Result, error) {
			user.Stats().WheelRolled++
			return e.Next()
		}),
		user.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*result.Result, error) {
			user.Stats().ItemsUsed++
			return e.Next()
		}),
	}
}

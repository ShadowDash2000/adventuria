package effects

import (
	"adventuria/internal/adventuria"
)

type DiceIncrementEffect struct {
	adventuria.EffectBase
}

func (ef *DiceIncrementEffect) Subscribe(callback adventuria.EffectCallback) {
	ef.PoolUnsubscribers(
		ef.User().OnBeforeRollMove().BindFunc(func(e *adventuria.OnBeforeRollMoveEvent) error {
			if i := ef.GetInt("value"); i != 0 {
				e.N += i
			}

			callback()

			return e.Next()
		}),
	)
}

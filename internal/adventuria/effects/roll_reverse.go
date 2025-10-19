package effects

import (
	"adventuria/internal/adventuria"
)

type RollReverseEffect struct {
	adventuria.EffectBase
}

func (ef *RollReverseEffect) Subscribe(callback adventuria.EffectCallback) {
	ef.PoolUnsubscribers(
		ef.User().OnBeforeRollMove().BindFunc(func(e *adventuria.OnBeforeRollMoveEvent) error {
			e.N *= -1

			callback()

			return e.Next()
		}),
	)
}

func (ef *RollReverseEffect) Verify(_ string) error {
	return nil
}

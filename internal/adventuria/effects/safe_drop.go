package effects

import (
	"adventuria/internal/adventuria"
)

type SafeDropEffect struct {
	adventuria.EffectBase
}

func (ef *SafeDropEffect) Subscribe(callback adventuria.EffectCallback) {
	ef.PoolUnsubscribers(
		ef.User().OnBeforeDrop().BindFunc(func(e *adventuria.OnBeforeDropEvent) error {
			e.IsSafeDrop = true

			callback()

			return e.Next()
		}),
	)
}

package effects

import (
	"adventuria/internal/adventuria"
)

type GiveWheelOnDoneEffect struct {
	adventuria.PersistentEffectBase
}

func (ef *GiveWheelOnDoneEffect) Subscribe(callback adventuria.EffectCallback) {
	ef.PoolUnsubscribers(
		ef.User().OnAfterDone().BindFunc(func(e *adventuria.OnAfterDoneEvent) error {
			ef.User().SetItemWheelsCount(ef.User().ItemWheelsCount() + 1)

			callback()

			return e.Next()
		}),
	)
}

package effects

import (
	"adventuria/internal/adventuria"
)

type JailEscapeEffect struct {
	adventuria.EffectBase
}

func (ef *JailEscapeEffect) Subscribe(callback adventuria.EffectCallback) {
	ef.PoolUnsubscribers(
		ef.User().OnAfterAction().BindFunc(func(e *adventuria.OnAfterActionEvent) error {
			ef.User().SetIsInJail(false)
			ef.User().SetDropsInARow(0)

			callback()

			return e.Next()
		}),
	)
}

func (ef *JailEscapeEffect) Verify(_ string) error {
	return nil
}

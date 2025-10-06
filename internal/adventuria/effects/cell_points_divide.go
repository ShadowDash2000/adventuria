package effects

import (
	"adventuria/internal/adventuria"
)

type CellPointsDivideEffect struct {
	adventuria.EffectBase
}

func (ef *CellPointsDivideEffect) Subscribe(callback adventuria.EffectCallback) {
	ef.PoolUnsubscribers(
		ef.User().OnBeforeDone().BindFunc(func(e *adventuria.OnBeforeDoneEvent) error {
			if i := ef.GetInt("value"); i != 0 {
				e.CellPointsDivide = i
			}

			callback()

			return e.Next()
		}),
	)
}

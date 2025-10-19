package effects

import (
	"adventuria/internal/adventuria"
	"strconv"
)

type DiceMultiplierEffect struct {
	adventuria.EffectBase
}

func (ef *DiceMultiplierEffect) Subscribe(callback adventuria.EffectCallback) {
	ef.PoolUnsubscribers(
		ef.User().OnBeforeRollMove().BindFunc(func(e *adventuria.OnBeforeRollMoveEvent) error {
			if i := ef.GetInt("value"); i != 0 {
				e.N *= i
			}

			callback()

			return e.Next()
		}),
	)
}

func (ef *DiceMultiplierEffect) Verify(value string) error {
	if _, err := strconv.Atoi(value); err != nil {
		return err
	}
	return nil
}

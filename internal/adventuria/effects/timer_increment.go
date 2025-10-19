package effects

import (
	"adventuria/internal/adventuria"
	"strconv"
)

type TimerIncrementEffect struct {
	adventuria.EffectBase
}

func (ef *TimerIncrementEffect) Subscribe(callback adventuria.EffectCallback) {
	ef.PoolUnsubscribers(
		ef.User().OnAfterAction().BindFunc(func(e *adventuria.OnAfterActionEvent) error {
			if i := ef.GetInt("value"); i != 0 {
				err := ef.User().Timer().AddSecondsTimeLimit(i)
				if err != nil {
					return err
				}
			}

			callback()

			return e.Next()
		}),
	)
}

func (ef *TimerIncrementEffect) Verify(value string) error {
	if _, err := strconv.Atoi(value); err != nil {
		return err
	}
	return nil
}

package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"strconv"
)

type TimerIncrementEffect struct {
	adventuria.EffectBase
}

func (ef *TimerIncrementEffect) Subscribe(
	user adventuria.User,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		user.OnAfterAction().BindFunc(func(e *adventuria.OnAfterActionEvent) error {
			if i := ef.GetInt("value"); i != 0 {
				err := user.Timer().AddSecondsTimeLimit(i)
				if err != nil {
					return err
				}
			}

			callback()

			return e.Next()
		}),
	}
}

func (ef *TimerIncrementEffect) Verify(value string) error {
	if _, err := strconv.Atoi(value); err != nil {
		return err
	}
	return nil
}

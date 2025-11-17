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
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *TimerIncrementEffect) DecodeValue(value string) (any, error) {
	return strconv.Atoi(value)
}

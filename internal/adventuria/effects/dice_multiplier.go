package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"strconv"
)

type DiceMultiplierEffect struct {
	adventuria.EffectBase
}

func (ef *DiceMultiplierEffect) Subscribe(
	user adventuria.User,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		user.OnBeforeRollMove().BindFunc(func(e *adventuria.OnBeforeRollMoveEvent) error {
			if i := ef.GetInt("value"); i != 0 {
				e.N *= i
			}

			callback()

			return e.Next()
		}),
	}
}

func (ef *DiceMultiplierEffect) Verify(value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *DiceMultiplierEffect) DecodeValue(value string) (any, error) {
	return strconv.Atoi(value)
}

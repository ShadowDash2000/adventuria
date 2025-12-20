package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"strconv"
)

type DiceIncrementEffect struct {
	adventuria.EffectBase
}

func (ef *DiceIncrementEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeRollMove().BindFunc(func(e *adventuria.OnBeforeRollMoveEvent) error {
			if i := ef.GetInt("value"); i != 0 {
				e.N += i

				callback()
			}

			return e.Next()
		}),
	}, nil
}

func (ef *DiceIncrementEffect) Verify(value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *DiceIncrementEffect) DecodeValue(value string) (any, error) {
	return strconv.Atoi(value)
}

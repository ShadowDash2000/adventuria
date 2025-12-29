package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"strconv"
)

type DiceMultiplierEffect struct {
	adventuria.EffectRecord
}

func (ef *DiceMultiplierEffect) CanUse(_ adventuria.EffectContext) bool {
	return true
}

func (ef *DiceMultiplierEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeRollMove().BindFunc(func(e *adventuria.OnBeforeRollMoveEvent) (*event.Result, error) {
			if i := ef.GetInt("value"); i != 0 {
				e.N *= i

				callback()
			}

			return e.Next()
		}),
	}, nil
}

func (ef *DiceMultiplierEffect) Verify(value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *DiceMultiplierEffect) DecodeValue(value string) (any, error) {
	return strconv.Atoi(value)
}

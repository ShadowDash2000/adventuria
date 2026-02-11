package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
	"strconv"
)

type DiceMultiplierEffect struct {
	adventuria.EffectRecord
}

func (ef *DiceMultiplierEffect) CanUse(_ adventuria.AppContext, _ adventuria.EffectContext) bool {
	return true
}

func (ef *DiceMultiplierEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeRollMove().BindFunc(func(e *adventuria.OnBeforeRollMoveEvent) (*result.Result, error) {
			if i := ef.GetInt("value"); i != 0 {
				e.N *= i

				callback(e.AppContext)
			}

			return e.Next()
		}),
	}, nil
}

func (ef *DiceMultiplierEffect) Verify(_ adventuria.AppContext, value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *DiceMultiplierEffect) DecodeValue(value string) (any, error) {
	return strconv.Atoi(value)
}

func (ef *DiceMultiplierEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

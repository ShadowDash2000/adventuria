package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"strconv"
)

type ChangeMaxGamePriceEffect struct {
	adventuria.EffectBase
}

func (ef *ChangeMaxGamePriceEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemAdd().BindFunc(func(e *adventuria.OnAfterItemAdd) error {
			if i := ef.GetInt("value"); i != 0 {
				ctx.User.LastAction().CustomGameFilter().MaxPrice = i

				callback()
			}

			return e.Next()
		}),
	}
}

func (ef *ChangeMaxGamePriceEffect) Verify(value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *ChangeMaxGamePriceEffect) DecodeValue(value string) (any, error) {
	return strconv.Atoi(value)
}

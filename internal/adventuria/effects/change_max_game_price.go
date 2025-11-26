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
	user adventuria.User,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		user.OnAfterItemAdd().BindFunc(func(e *adventuria.OnAfterItemAdd) error {
			if i := ef.GetInt("value"); i != 0 {
				user.LastAction().CustomGameFilter().MaxPrice = i
			}

			callback()

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

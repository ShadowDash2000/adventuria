package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type NothingEffect struct {
	adventuria.EffectBase
}

func (ef *NothingEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemAdd().BindFunc(func(e *adventuria.OnAfterItemAdd) error {
			callback()

			return e.Next()
		}),
	}
}

func (ef *NothingEffect) Verify(_ string) error {
	return nil
}

func (ef *NothingEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}

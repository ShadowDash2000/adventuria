package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type SafeDropEffect struct {
	adventuria.EffectBase
}

func (ef *SafeDropEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		ctx.User.OnBeforeDrop().BindFunc(func(e *adventuria.OnBeforeDropEvent) error {
			e.IsSafeDrop = true

			callback()

			return e.Next()
		}),
	}
}

func (ef *SafeDropEffect) Verify(_ string) error {
	return nil
}

func (ef *SafeDropEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}

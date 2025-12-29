package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type SafeDropEffect struct {
	adventuria.EffectRecord
}

func (ef *SafeDropEffect) CanUse(_ adventuria.EffectContext) bool {
	return true
}

func (ef *SafeDropEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeDrop().BindFunc(func(e *adventuria.OnBeforeDropEvent) (*event.Result, error) {
			e.IsSafeDrop = true

			callback()

			return e.Next()
		}),
	}, nil
}

func (ef *SafeDropEffect) Verify(_ string) error {
	return nil
}

func (ef *SafeDropEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}

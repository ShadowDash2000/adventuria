package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type RollReverseEffect struct {
	adventuria.EffectBase
}

func (ef *RollReverseEffect) Subscribe(
	user adventuria.User,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		user.OnBeforeRollMove().BindFunc(func(e *adventuria.OnBeforeRollMoveEvent) error {
			e.N *= -1

			callback()

			return e.Next()
		}),
	}
}

func (ef *RollReverseEffect) Verify(_ string) error {
	return nil
}

func (ef *RollReverseEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}

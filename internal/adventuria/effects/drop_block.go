package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type DropBlockedEffect struct {
	adventuria.EffectBase
}

func (ef *DropBlockedEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeDrop().BindFunc(func(e *adventuria.OnBeforeDropEvent) (*event.Result, error) {
			e.IsDropBlocked = true
			callback()
			return e.Next()
		}),
	}, nil
}

func (ef *DropBlockedEffect) Verify(_ string) error {
	return nil
}

func (ef *DropBlockedEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}

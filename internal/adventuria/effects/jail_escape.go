package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type JailEscapeEffect struct {
	adventuria.EffectBase
}

func (ef *JailEscapeEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		ctx.User.OnAfterAction().BindFunc(func(e *adventuria.OnAfterActionEvent) error {
			ctx.User.SetIsInJail(false)
			ctx.User.SetDropsInARow(0)

			callback()

			return e.Next()
		}),
	}
}

func (ef *JailEscapeEffect) Verify(_ string) error {
	return nil
}

func (ef *JailEscapeEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}

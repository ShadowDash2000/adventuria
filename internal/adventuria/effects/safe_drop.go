package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
	"errors"
)

type SafeDropEffect struct {
	adventuria.EffectRecord
}

func (ef *SafeDropEffect) CanUse(_ adventuria.AppContext, _ adventuria.EffectContext) bool {
	return true
}

func (ef *SafeDropEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeDrop().BindFunc(func(e *adventuria.OnBeforeDropEvent) (*result.Result, error) {
			cell, ok := ctx.User.CurrentCell()
			if !ok {
				return result.Err("internal error: current cell not found"),
					errors.New("safeDrop: current cell not found")
			}

			if cell.IsSafeDrop() {
				return e.Next()
			}

			e.IsSafeDrop = true

			callback(e.AppContext)

			return e.Next()
		}),
	}, nil
}

func (ef *SafeDropEffect) Verify(_ adventuria.AppContext, _ string) error {
	return nil
}

func (ef *SafeDropEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

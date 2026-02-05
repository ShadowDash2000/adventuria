package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type RollReverseEffect struct {
	adventuria.EffectRecord
}

func (ef *RollReverseEffect) CanUse(_ adventuria.AppContext, _ adventuria.EffectContext) bool {
	return true
}

func (ef *RollReverseEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeRollMove().BindFunc(func(e *adventuria.OnBeforeRollMoveEvent) (*event.Result, error) {
			e.N *= -1

			callback(e.AppContext)

			return e.Next()
		}),
	}, nil
}

func (ef *RollReverseEffect) Verify(_ adventuria.AppContext, _ string) error {
	return nil
}

func (ef *RollReverseEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

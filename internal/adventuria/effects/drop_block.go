package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type DropBlockedEffect struct {
	adventuria.EffectRecord
}

func (ef *DropBlockedEffect) CanUse(_ adventuria.AppContext, _ adventuria.EffectContext) bool {
	return true
}

func (ef *DropBlockedEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeDropCheck().BindFunc(func(e *adventuria.OnBeforeDropCheckEvent) (*event.Result, error) {
			e.IsDropBlocked = true
			return e.Next()
		}),
		ctx.User.OnAfterDone().BindFunc(func(e *adventuria.OnAfterDoneEvent) (*event.Result, error) {
			callback(e.AppContext)
			return e.Next()
		}),
	}, nil
}

func (ef *DropBlockedEffect) Verify(_ adventuria.AppContext, _ string) error {
	return nil
}

func (ef *DropBlockedEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

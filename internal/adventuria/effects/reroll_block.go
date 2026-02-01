package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type RerollBlockedEffect struct {
	adventuria.EffectRecord
}

func (ef *RerollBlockedEffect) CanUse(_ adventuria.EffectContext) bool {
	return true
}

func (ef *RerollBlockedEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeRerollCheck().BindFunc(func(e *adventuria.OnBeforeRerollCheckEvent) (*event.Result, error) {
			e.IsRerollBlocked = true
			return e.Next()
		}),
		ctx.User.OnAfterDone().BindFunc(func(e *adventuria.OnAfterDoneEvent) (*event.Result, error) {
			callback()
			return e.Next()
		}),
	}, nil
}

func (ef *RerollBlockedEffect) Verify(_ string) error {
	return nil
}

func (ef *RerollBlockedEffect) GetVariants(_ adventuria.EffectContext) any {
	return nil
}

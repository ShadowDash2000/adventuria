package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type NoCoinsForDoneEffect struct {
	adventuria.EffectRecord
}

func (ef *NoCoinsForDoneEffect) CanUse(_ adventuria.AppContext, _ adventuria.EffectContext) bool {
	return true
}

func (ef *NoCoinsForDoneEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeDone().BindFunc(func(e *adventuria.OnBeforeDoneEvent) (*event.Result, error) {
			e.CellCoins = 0
			callback(e.AppContext)
			return e.Next()
		}),
	}, nil
}

func (ef *NoCoinsForDoneEffect) Verify(_ adventuria.AppContext, _ string) error {
	return nil
}

func (ef *NoCoinsForDoneEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

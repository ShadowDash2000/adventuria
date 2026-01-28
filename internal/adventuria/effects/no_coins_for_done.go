package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type NoCoinsForDoneEffect struct {
	adventuria.EffectRecord
}

func (ef *NoCoinsForDoneEffect) CanUse(_ adventuria.EffectContext) bool {
	return true
}

func (ef *NoCoinsForDoneEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeDone().BindFunc(func(e *adventuria.OnBeforeDoneEvent) (*event.Result, error) {
			e.CellCoins = 0
			callback()
			return e.Next()
		}),
	}, nil
}

func (ef *NoCoinsForDoneEffect) Verify(_ string) error {
	return nil
}

func (ef *NoCoinsForDoneEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}

func (ef *NoCoinsForDoneEffect) GetVariants(ctx adventuria.EffectContext) any {
	return nil
}

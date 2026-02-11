package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
)

type DebuffBlockEffect struct {
	adventuria.EffectRecord
}

func (ef *DebuffBlockEffect) CanUse(_ adventuria.AppContext, _ adventuria.EffectContext) bool {
	return true
}

func (ef *DebuffBlockEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeItemAdd().BindFunc(func(e *adventuria.OnBeforeItemAdd) (*result.Result, error) {
			if e.ItemRecord.Type() == "debuff" {
				e.ShouldAddItem = false
				callback(e.AppContext)
			}
			return e.Next()
		}),
	}, nil
}

func (ef *DebuffBlockEffect) Verify(_ adventuria.AppContext, _ string) error {
	return nil
}

func (ef *DebuffBlockEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

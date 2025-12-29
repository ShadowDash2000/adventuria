package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type DebuffBlockEffect struct {
	adventuria.EffectRecord
}

func (ef *DebuffBlockEffect) CanUse(_ adventuria.EffectContext) bool {
	return true
}

func (ef *DebuffBlockEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeItemAdd().BindFunc(func(e *adventuria.OnBeforeItemAdd) (*event.Result, error) {
			if e.ItemRecord.Type() == "debuff" {
				e.ShouldAddItem = false
				callback()
			}
			return e.Next()
		}),
	}, nil
}

func (ef *DebuffBlockEffect) Verify(_ string) error {
	return nil
}

func (ef *DebuffBlockEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}

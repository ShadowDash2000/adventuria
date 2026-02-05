package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"strconv"
)

type ReplaceDiceRollEffect struct {
	adventuria.EffectRecord
}

func (ef *ReplaceDiceRollEffect) CanUse(_ adventuria.AppContext, _ adventuria.EffectContext) bool {
	return true
}

func (ef *ReplaceDiceRollEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeRollMove().BindFunc(func(e *adventuria.OnBeforeRollMoveEvent) (*event.Result, error) {
			e.N = ef.GetInt("value")

			callback(e.AppContext)

			return e.Next()
		}),
	}, nil
}

func (ef *ReplaceDiceRollEffect) Verify(_ adventuria.AppContext, value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *ReplaceDiceRollEffect) DecodeValue(value string) (any, error) {
	return strconv.Atoi(value)
}

func (ef *ReplaceDiceRollEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

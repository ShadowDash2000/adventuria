package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"strconv"
)

type ReplaceDiceRollEffect struct {
	adventuria.EffectBase
}

func (ef *ReplaceDiceRollEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		ctx.User.OnBeforeRollMove().BindFunc(func(e *adventuria.OnBeforeRollMoveEvent) error {
			e.N = ef.GetInt("value")

			callback()

			return e.Next()
		}),
	}
}

func (ef *ReplaceDiceRollEffect) Verify(value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *ReplaceDiceRollEffect) DecodeValue(value string) (any, error) {
	return strconv.Atoi(value)
}

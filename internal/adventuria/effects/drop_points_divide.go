package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"strconv"
)

type DropPointsDivideEffect struct {
	adventuria.EffectBase
}

func (ef *DropPointsDivideEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeDrop().BindFunc(func(e *adventuria.OnBeforeDropEvent) error {
			if i := ef.GetInt("value"); i != 0 {
				e.PointsForDrop = e.PointsForDrop / i

				callback()
			}

			return e.Next()
		}),
	}, nil
}

func (ef *DropPointsDivideEffect) Verify(value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *DropPointsDivideEffect) DecodeValue(value string) (any, error) {
	return strconv.Atoi(value)
}

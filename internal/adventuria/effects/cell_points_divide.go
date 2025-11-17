package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"strconv"
)

type CellPointsDivideEffect struct {
	adventuria.EffectBase
}

func (ef *CellPointsDivideEffect) Subscribe(
	user adventuria.User,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		user.OnBeforeDone().BindFunc(func(e *adventuria.OnBeforeDoneEvent) error {
			if i := ef.GetInt("value"); i != 0 {
				e.CellPointsDivide = i
			}

			callback()

			return e.Next()
		}),
	}
}

func (ef *CellPointsDivideEffect) Verify(value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *CellPointsDivideEffect) DecodeValue(value string) (any, error) {
	return strconv.Atoi(value)
}

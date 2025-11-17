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
	if _, err := strconv.Atoi(value); err != nil {
		return err
	}
	return nil
}

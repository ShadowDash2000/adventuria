package effects

import (
	"adventuria/internal/adventuria"
	"strconv"
)

type PointsIncrementEffect struct {
	adventuria.EffectBase
}

func (ef *PointsIncrementEffect) Subscribe(callback adventuria.EffectCallback) {
	ef.PoolUnsubscribers(
		ef.User().OnAfterAction().BindFunc(func(e *adventuria.OnAfterActionEvent) error {
			if i := ef.GetInt("value"); i != 0 {
				ef.User().SetPoints(ef.User().Points() + i)
			}

			callback()

			return e.Next()
		}),
	)
}

func (ef *PointsIncrementEffect) Verify(value string) error {
	if _, err := strconv.Atoi(value); err != nil {
		return err
	}
	return nil
}

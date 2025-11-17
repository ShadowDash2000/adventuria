package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"strconv"
)

type PointsIncrementEffect struct {
	adventuria.EffectBase
}

func (ef *PointsIncrementEffect) Subscribe(
	user adventuria.User,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		user.OnAfterAction().BindFunc(func(e *adventuria.OnAfterActionEvent) error {
			if i := ef.GetInt("value"); i != 0 {
				user.SetPoints(user.Points() + i)
			}

			callback()

			return e.Next()
		}),
	}
}

func (ef *PointsIncrementEffect) Verify(value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *PointsIncrementEffect) DecodeValue(value string) (any, error) {
	return strconv.Atoi(value)
}

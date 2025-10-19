package effects

import (
	"adventuria/internal/adventuria"
)

type GiveWheelOnNewLapEffect struct {
	adventuria.PersistentEffectBase
}

func (ef *GiveWheelOnNewLapEffect) Subscribe() {
	ef.PoolUnsubscribers(
		ef.User().OnNewLap().BindFunc(func(e *adventuria.OnNewLapEvent) error {
			ef.User().SetItemWheelsCount(ef.User().ItemWheelsCount() + e.Laps)

			return e.Next()
		}),
	)
}

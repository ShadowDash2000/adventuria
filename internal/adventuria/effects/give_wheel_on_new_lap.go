package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
)

type GiveWheelOnNewLapEffect struct{}

func (ef *GiveWheelOnNewLapEffect) Subscribe(user adventuria.User) []event.Unsubscribe {
	return []event.Unsubscribe{
		user.OnNewLap().BindFunc(func(e *adventuria.OnNewLapEvent) (*result.Result, error) {
			user.SetItemWheelsCount(user.ItemWheelsCount() + e.Laps)

			return e.Next()
		}),
	}
}

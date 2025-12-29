package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type GiveWheelOnNewLapEffect struct{}

func (ef *GiveWheelOnNewLapEffect) Subscribe(user adventuria.User) []event.Unsubscribe {
	return []event.Unsubscribe{
		user.OnNewLap().BindFunc(func(e *adventuria.OnNewLapEvent) (*event.Result, error) {
			user.SetItemWheelsCount(user.ItemWheelsCount() + e.Laps)

			return e.Next()
		}),
	}
}

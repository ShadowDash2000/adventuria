package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
)

type GiveWheelOnNewLapEffect struct{}

func (ef *GiveWheelOnNewLapEffect) Subscribe(player adventuria.Player) []event.Unsubscribe {
	return []event.Unsubscribe{
		player.OnNewLap().BindFunc(func(e *adventuria.OnNewLapEvent) (*result.Result, error) {
			player.Progress().AddItemWheelsCount(e.Laps)

			return e.Next()
		}),
	}
}

package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
)

type GiveWheelOnDoneEffect struct{}

func (ef *GiveWheelOnDoneEffect) Subscribe(player adventuria.Player) []event.Unsubscribe {
	return []event.Unsubscribe{
		player.OnAfterDone().BindFunc(func(e *adventuria.OnAfterDoneEvent) (*result.Result, error) {
			player.Progress().AddItemWheelsCount(1)

			return e.Next()
		}),
	}
}

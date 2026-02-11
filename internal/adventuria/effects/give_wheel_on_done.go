package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
)

type GiveWheelOnDoneEffect struct{}

func (ef *GiveWheelOnDoneEffect) Subscribe(user adventuria.User) []event.Unsubscribe {
	return []event.Unsubscribe{
		user.OnAfterDone().BindFunc(func(e *adventuria.OnAfterDoneEvent) (*result.Result, error) {
			user.SetItemWheelsCount(user.ItemWheelsCount() + 1)

			return e.Next()
		}),
	}
}

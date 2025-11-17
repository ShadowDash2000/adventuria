package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type GiveWheelOnDoneEffect struct {
	adventuria.PersistentEffectBase
}

func (ef *GiveWheelOnDoneEffect) Subscribe(user adventuria.User) []event.Unsubscribe {
	return []event.Unsubscribe{
		user.OnAfterDone().BindFunc(func(e *adventuria.OnAfterDoneEvent) error {
			user.SetItemWheelsCount(user.ItemWheelsCount() + 1)

			return e.Next()
		}),
	}
}

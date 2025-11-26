package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type NoTimeLimitEffect struct {
	adventuria.EffectBase
}

func (ef *NoTimeLimitEffect) Subscribe(
	user adventuria.User,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		user.OnAfterItemAdd().BindFunc(func(e *adventuria.OnAfterItemAdd) error {
			user.LastAction().CustomGameFilter().MinCampaignTime = -1
			user.LastAction().CustomGameFilter().MaxCampaignTime = -1

			callback()

			return e.Next()
		}),
	}
}

func (ef *NoTimeLimitEffect) Verify(_ string) error {
	return nil
}

func (ef *NoTimeLimitEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}

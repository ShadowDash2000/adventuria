package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type NoTimeLimitEffect struct {
	adventuria.EffectBase
}

func (ef *NoTimeLimitEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemAdd().BindFunc(func(e *adventuria.OnAfterItemAdd) error {
			ctx.User.LastAction().CustomGameFilter().MinCampaignTime = -1
			ctx.User.LastAction().CustomGameFilter().MaxCampaignTime = -1

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

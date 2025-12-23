package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type JailEscapeEffect struct {
	adventuria.EffectBase
}

func (ef *JailEscapeEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*event.Result, error) {
			if e.InvItemId == ctx.InvItemID {
				if !ctx.User.IsInJail() {
					return &event.Result{
						Success: false,
						Error:   "user isn't in jail",
					}, nil
				}

				ctx.User.SetIsInJail(false)
				ctx.User.SetDropsInARow(0)

				callback()
			}

			return e.Next()
		}),
	}, nil
}

func (ef *JailEscapeEffect) Verify(_ string) error {
	return nil
}

func (ef *JailEscapeEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}

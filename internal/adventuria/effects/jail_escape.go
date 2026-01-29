package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type JailEscapeEffect struct {
	adventuria.EffectRecord
}

func (ef *JailEscapeEffect) CanUse(ctx adventuria.EffectContext) bool {
	return ctx.User.IsInJail()
}

func (ef *JailEscapeEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*event.Result, error) {
			if e.InvItemId == ctx.InvItemID {
				ctx.User.SetIsInJail(false)
				ctx.User.SetDropsInARow(0)
				ctx.User.LastAction().SetCanMove(true)

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

func (ef *JailEscapeEffect) GetVariants(ctx adventuria.EffectContext) any {
	return nil
}

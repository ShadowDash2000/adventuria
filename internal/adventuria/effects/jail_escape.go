package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
)

type JailEscapeEffect struct {
	adventuria.EffectRecord
}

func (ef *JailEscapeEffect) CanUse(_ adventuria.AppContext, ctx adventuria.EffectContext) bool {
	return ctx.Player.Progress().IsInJail()
}

func (ef *JailEscapeEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.Player.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*result.Result, error) {
			if e.InvItemId == ctx.InvItemID {
				ctx.Player.Progress().SetIsInJail(false)
				ctx.Player.Progress().SetDropsInARow(0)
				ctx.Player.LastAction().SetCanMove(true)

				callback(e.AppContext)
			}

			return e.Next()
		}),
	}, nil
}

func (ef *JailEscapeEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

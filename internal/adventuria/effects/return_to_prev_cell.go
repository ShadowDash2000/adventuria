package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
	"fmt"
)

type ReturnToPrevCellEffect struct {
	adventuria.EffectRecord
}

func (ef *ReturnToPrevCellEffect) CanUse(appCtx adventuria.AppContext, ctx adventuria.EffectContext) bool {
	latestDiceRoll := ctx.Player.LastAction().CellsPassed()
	if latestDiceRoll == 0 {
		return false
	}

	return !adventuria.GameActions.CanDo(appCtx, ctx.Player, "done")
}

func (ef *ReturnToPrevCellEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.Player.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*result.Result, error) {
			if e.InvItemId != ctx.InvItemID {
				return e.Next()
			}

			latestDiceRoll := ctx.Player.LastAction().CellsPassed()
			_, err := ctx.Player.Move(e.AppContext, -latestDiceRoll)
			if err != nil {
				return result.Err("internal error: failed to move to the previous cell"),
					fmt.Errorf("returnToPrevCell: %w", err)
			}

			ctx.Player.LastAction().SetCanMove(true)

			callback(e.AppContext)

			return e.Next()
		}),
	}, nil
}

func (ef *ReturnToPrevCellEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

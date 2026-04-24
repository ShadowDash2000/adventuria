package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
	"errors"
)

type SaveFromHoleEffect struct {
	adventuria.EffectRecord
}

func (ef *SaveFromHoleEffect) CanUse(_ adventuria.AppContext, _ adventuria.EffectContext) bool {
	return true
}

func (ef *SaveFromHoleEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.Player.OnBeforeTeleportOnCell().BindFunc(func(e *adventuria.OnBeforeTeleportOnCell) (*result.Result, error) {
			if e.SkipTeleport {
				return e.Next()
			}

			targetCellOrder, ok := adventuria.GameCells.GetGlobalOrderById(e.CellId)
			if !ok {
				return result.Err("internal error: target cell not found"),
					errors.New("saveFromHole: current cell not found")
			}

			if targetCellOrder < ctx.Player.Progress().GlobalCurrentCellOrder() {
				e.SkipTeleport = true
				callback(e.AppContext)
			}

			return e.Next()
		}),
	}, nil
}

func (ef *SaveFromHoleEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
	"math/rand"
)

type TeleportToRandomCellEffect struct {
	adventuria.EffectRecord
}

func (ef *TeleportToRandomCellEffect) CanUse(appCtx adventuria.AppContext, ctx adventuria.EffectContext) bool {
	if adventuria.GameActions.CanDo(appCtx, ctx.User, "rollDice") {
		return true
	}

	if adventuria.GameActions.HasActionsInCategories(appCtx, ctx.User, []string{"wheel_roll", "on_cell"}) {
		return false
	}

	currentCell, ok := ctx.User.CurrentCell()
	if !ok {
		return false
	}

	canDone := adventuria.GameActions.CanDo(appCtx, ctx.User, "done")
	canDrop := adventuria.GameActions.CanDo(appCtx, ctx.User, "drop")

	if canDone && !canDrop {
		if currentCell.Type() != "jail" {
			return false
		}
	}

	return true
}

func (ef *TeleportToRandomCellEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*result.Result, error) {
			if ctx.InvItemID == e.InvItemId {
				totalCells := adventuria.GameCells.Count()
				if totalCells <= 1 {
					return e.Next()
				}

				randomCellOrder := rand.Intn(totalCells - 1)
				if randomCellOrder >= ctx.User.CurrentCellOrder() {
					randomCellOrder++
				}

				_, err := ctx.User.MoveToCellOrder(e.AppContext, randomCellOrder)
				if err != nil {
					return result.Err("internal error: failed to move to the cell"), err
				}

				callback(e.AppContext)
			}

			return e.Next()
		}),
	}, nil
}

func (ef *TeleportToRandomCellEffect) Verify(_ adventuria.AppContext, _ string) error {
	return nil
}

func (ef *TeleportToRandomCellEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

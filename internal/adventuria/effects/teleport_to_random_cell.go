package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
	"errors"
	"math/rand"
)

type TeleportToRandomCellEffect struct {
	adventuria.EffectRecord
}

func (ef *TeleportToRandomCellEffect) CanUse(appCtx adventuria.AppContext, ctx adventuria.EffectContext) bool {
	if adventuria.GameActions.CanDo(appCtx, ctx.Player, "rollDice") {
		return true
	}

	if adventuria.GameActions.HasActionsInCategories(appCtx, ctx.Player, []string{"wheel_roll", "on_cell"}) {
		return false
	}

	currentCell, ok := ctx.Player.Progress().CurrentCell()
	if !ok {
		return false
	}

	canDone := adventuria.GameActions.CanDo(appCtx, ctx.Player, "done")
	canDrop := adventuria.GameActions.CanDo(appCtx, ctx.Player, "drop")

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
		ctx.Player.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*result.Result, error) {
			if ctx.InvItemID == e.InvItemId {
				totalCells := adventuria.GameCells.CountGlobal()
				if totalCells <= 1 {
					return e.Next()
				}

				randomCellOrder := rand.Intn(totalCells - 1)
				if randomCellOrder >= ctx.Player.Progress().GlobalCurrentCellOrder() {
					randomCellOrder++
				}

				randomCell, ok := adventuria.GameCells.GetByGlobalOrder(randomCellOrder)
				if !ok {
					return result.Err("internal error: random cell not found"),
						errors.New("teleportToRandomCell: random cell not found")
				}

				_, err := ctx.Player.MoveToCellId(e.AppContext, randomCell.ID())
				if err != nil {
					return result.Err("internal error: failed to move to the cell"), err
				}

				callback(e.AppContext)
			}

			return e.Next()
		}),
	}, nil
}

func (ef *TeleportToRandomCellEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

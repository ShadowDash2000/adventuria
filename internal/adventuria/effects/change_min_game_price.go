package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"adventuria/pkg/event"
	"fmt"
	"strconv"
)

type ChangeMinGamePriceEffect struct {
	adventuria.EffectBase
}

func (ef *ChangeMinGamePriceEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*event.Result, error) {
			if e.InvItemId == ctx.InvItemID {
				if ok := adventuria.GameActions.CanDo(ctx.User, "rollWheel"); !ok {
					return &event.Result{
						Success: false,
						Error:   "user can't perform rollWheel action",
					}, nil
				}

				cell, ok := ctx.User.CurrentCell()
				if !ok {
					return &event.Result{
						Success: false,
						Error:   "current cell not found",
					}, nil
				}

				cellGame, ok := cell.(*cells.CellGame)
				if !ok {
					return &event.Result{
						Success: false,
						Error:   "current cell isn't game cell",
					}, nil
				}

				if cell.Type() != "game" {
					return &event.Result{
						Success: false,
						Error:   "current cell isn't game cell",
					}, nil
				}

				if i := ef.GetInt("value"); i != 0 {
					ctx.User.LastAction().CustomActivityFilter().MinPrice = i
					if err := cellGame.RefreshItems(ctx.User); err != nil {
						return &event.Result{
							Success: false,
							Error:   "internal error: can't refresh cell items in \"change_min_game_price\" effect",
						}, fmt.Errorf("changeMinGamePrice: %w", err)
					}

					callback()
				}
			}

			return e.Next()
		}),
	}, nil
}

func (ef *ChangeMinGamePriceEffect) Verify(value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *ChangeMinGamePriceEffect) DecodeValue(value string) (any, error) {
	return strconv.Atoi(value)
}

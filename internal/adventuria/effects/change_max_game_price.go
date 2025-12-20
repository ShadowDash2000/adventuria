package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"adventuria/pkg/event"
	"errors"
	"fmt"
	"strconv"
)

type ChangeMaxGamePriceEffect struct {
	adventuria.EffectBase
}

func (ef *ChangeMaxGamePriceEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) error {
			if e.InvItemId == ctx.InvItemID {
				if ok := adventuria.GameActions.CanDo(ctx.User, "rollWheel"); !ok {
					return errors.New("noTimeLimit: user can't do rollWheel action")
				}

				cell, ok := ctx.User.CurrentCell()
				if !ok {
					return errors.New("changeMinGamePrice: current cell not found")
				}

				cellGame, ok := cell.(*cells.CellGame)
				if !ok {
					return errors.New("changeMinGamePrice: current cell isn't game cell")
				}

				if i := ef.GetInt("value"); i != 0 {
					ctx.User.LastAction().CustomGameFilter().MaxPrice = i
					if err := cellGame.CheckCustomFilter(ctx.User); err != nil {
						return fmt.Errorf("changeMinGamePrice: %w", err)
					}

					callback()
				}
			}

			return e.Next()
		}),
	}, nil
}

func (ef *ChangeMaxGamePriceEffect) Verify(value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *ChangeMaxGamePriceEffect) DecodeValue(value string) (any, error) {
	return strconv.Atoi(value)
}

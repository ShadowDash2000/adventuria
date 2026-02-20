package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
	"errors"
	"fmt"
)

type StayOnCellAfterDoneEffect struct {
	adventuria.EffectRecord
}

func (ef *StayOnCellAfterDoneEffect) CanUse(_ adventuria.AppContext, _ adventuria.EffectContext) bool {
	return true
}

func (ef *StayOnCellAfterDoneEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterDone().BindFunc(func(e *adventuria.OnAfterDoneEvent) (*result.Result, error) {
			lastAction := ctx.User.LastAction()
			if lastAction.Type() != "done" {
				return e.Next()
			}

			cell, ok := ctx.User.CurrentCell()
			if !ok {
				return result.Err("internal error: current cell not found"),
					errors.New("stayOnCellAfterDone: current cell not found")
			}

			cellWheel, ok := cell.(adventuria.CellWheel)
			if !ok {
				return result.Err("current cell isn't wheel cell"), nil
			}

			err := e.App.Save(lastAction.ProxyRecord())
			if err != nil {
				return result.Err("internal error: failed to save lastest action"),
					fmt.Errorf("stayOnCellAfterDone: %w", err)
			}

			err = cellWheel.RefreshItems(e.AppContext, ctx.User)
			if err != nil {
				return result.Err("internal error: failed to refresh items"),
					fmt.Errorf("stayOnCellAfterDone: %w", err)
			}

			lastAction.MarkAsNew()
			lastAction.SetCanMove(false)
			lastAction.SetType("rollDice")
			lastAction.ClearCustomActivityFilter()

			callback(e.AppContext)

			return e.Next()
		}),
	}, nil
}

func (ef *StayOnCellAfterDoneEffect) Verify(_ adventuria.AppContext, _ string) error {
	return nil
}

func (ef *StayOnCellAfterDoneEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

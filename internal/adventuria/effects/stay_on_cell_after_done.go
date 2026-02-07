package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
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
		ctx.User.OnAfterDone().BindFunc(func(e *adventuria.OnAfterDoneEvent) (*event.Result, error) {
			cell, ok := ctx.User.CurrentCell()
			if !ok {
				return &event.Result{
					Success: false,
					Error:   "current cell not found",
				}, nil
			}

			// we can "done" only on wheel cell, so we won't need unnesseccary checks
			cellWheel, ok := cell.(adventuria.CellWheel)
			if !ok {
				return &event.Result{
					Success: false,
					Error:   "current cell isn't wheel cell",
				}, nil
			}

			lastAction := ctx.User.LastAction()
			err := e.App.Save(lastAction.ProxyRecord())
			if err != nil {
				return &event.Result{
					Success: false,
					Error:   "internal error: failed to save lastest action",
				}, fmt.Errorf("stayOnCellAfterDone: %w", err)
			}

			err = cellWheel.RefreshItems(e.AppContext, ctx.User)
			if err != nil {
				return &event.Result{
					Success: false,
					Error:   "internal error: failed to refresh items",
				}, fmt.Errorf("stayOnCellAfterDone: %w", err)
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

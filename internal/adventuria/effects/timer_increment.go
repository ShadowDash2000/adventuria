package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"fmt"
	"strconv"
)

type TimerIncrementEffect struct {
	adventuria.EffectRecord
}

func (ef *TimerIncrementEffect) CanUse(_ adventuria.EffectContext) bool {
	return true
}

func (ef *TimerIncrementEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterAction().BindFunc(func(e *adventuria.OnAfterActionEvent) (*event.Result, error) {
			if i := ef.GetInt("value"); i != 0 {
				err := ctx.User.Timer().AddSecondsTimeLimit(i)
				if err != nil {
					return &event.Result{
						Success: false,
						Error:   "internal error: failed to increment timer",
					}, fmt.Errorf("timerIncrementEffect: %w", err)
				}

				callback()
			}

			return e.Next()
		}),
	}, nil
}

func (ef *TimerIncrementEffect) Verify(value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *TimerIncrementEffect) DecodeValue(value string) (any, error) {
	return strconv.Atoi(value)
}

func (ef *TimerIncrementEffect) GetVariants(ctx adventuria.EffectContext) any {
	return nil
}

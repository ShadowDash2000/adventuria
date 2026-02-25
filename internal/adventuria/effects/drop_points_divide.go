package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
	"errors"
	"strconv"
)

type DropPointsDivideEffect struct {
	adventuria.EffectRecord
}

func (ef *DropPointsDivideEffect) CanUse(_ adventuria.AppContext, _ adventuria.EffectContext) bool {
	return true
}

func (ef *DropPointsDivideEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeDrop().BindFunc(func(e *adventuria.OnBeforeDropEvent) (*result.Result, error) {
			cell, ok := ctx.User.CurrentCell()
			if !ok {
				return result.Err("internal error: current cell not found"),
					errors.New("dropPointsDivide: current cell not found")
			}

			if cell.IsSafeDrop() {
				return e.Next()
			}

			if i := ef.GetInt("value"); i != 0 {
				e.PointsForDrop = e.PointsForDrop / i

				callback(e.AppContext)
			}

			return e.Next()
		}),
	}, nil
}

func (ef *DropPointsDivideEffect) Verify(_ adventuria.AppContext, value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *DropPointsDivideEffect) DecodeValue(value string) (any, error) {
	return strconv.Atoi(value)
}

func (ef *DropPointsDivideEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

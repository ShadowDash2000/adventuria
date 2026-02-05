package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"strconv"
)

type CellPointsDivideEffect struct {
	adventuria.EffectRecord
}

func (ef *CellPointsDivideEffect) CanUse(_ adventuria.AppContext, _ adventuria.EffectContext) bool {
	return true
}

func (ef *CellPointsDivideEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeDone().BindFunc(func(e *adventuria.OnBeforeDoneEvent) (*event.Result, error) {
			if i := ef.GetInt("value"); i != 0 {
				e.CellPoints = e.CellPoints / i

				callback(e.AppContext)
			}

			return e.Next()
		}),
		ctx.User.OnAfterMove().BindFunc(func(e *adventuria.OnAfterMoveEvent) (*event.Result, error) {
			callback(e.AppContext)
			return e.Next()
		}),
	}, nil
}

func (ef *CellPointsDivideEffect) Verify(_ adventuria.AppContext, value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *CellPointsDivideEffect) DecodeValue(value string) (any, error) {
	return strconv.Atoi(value)
}

func (ef *CellPointsDivideEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

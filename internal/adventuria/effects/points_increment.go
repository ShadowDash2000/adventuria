package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"strconv"
)

type PointsIncrementEffect struct {
	adventuria.EffectRecord
}

func (ef *PointsIncrementEffect) CanUse(_ adventuria.EffectContext) bool {
	return true
}

func (ef *PointsIncrementEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterAction().BindFunc(func(e *adventuria.OnAfterActionEvent) (*event.Result, error) {
			if i := ef.GetInt("value"); i != 0 {
				ctx.User.SetPoints(ctx.User.Points() + i)

				callback()
			}

			return e.Next()
		}),
	}, nil
}

func (ef *PointsIncrementEffect) Verify(value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *PointsIncrementEffect) DecodeValue(value string) (any, error) {
	return strconv.Atoi(value)
}

func (ef *PointsIncrementEffect) GetVariants(ctx adventuria.EffectContext) any {
	return nil
}

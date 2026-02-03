package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"errors"
	"strconv"
)

type DiscountPriceDivideEffect struct {
	adventuria.EffectRecord
}

func (ef *DiscountPriceDivideEffect) CanUse(_ adventuria.EffectContext) bool {
	return true
}

func (ef *DiscountPriceDivideEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeItemBuy().BindFunc(func(e *adventuria.OnBeforeItemBuy) (*event.Result, error) {
			e.Price /= ef.GetInt("value")
			callback()
			return e.Next()
		}),
		ctx.User.OnBuyGetVariants().BindFunc(func(e *adventuria.OnBuyGetVariants) (*event.Result, error) {
			e.Price /= ef.GetInt("value")
			return e.Next()
		}),
	}, nil
}

func (ef *DiscountPriceDivideEffect) Verify(value string) error {
	num, err := ef.DecodeValue(value)
	if num == 0 {
		return errors.New("discountPriceDivide: value must not be 0")
	}
	return err
}

func (ef *DiscountPriceDivideEffect) DecodeValue(value string) (int, error) {
	return strconv.Atoi(value)
}

func (ef *DiscountPriceDivideEffect) GetVariants(_ adventuria.EffectContext) any {
	return nil
}

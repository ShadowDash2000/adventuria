package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"errors"
	"fmt"
	"strings"
)

type ChangeDicesEffect struct {
	adventuria.EffectRecord
}

func (ef *ChangeDicesEffect) CanUse(_ adventuria.AppContext, _ adventuria.EffectContext) bool {
	return true
}

func (ef *ChangeDicesEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnBeforeRoll().BindFunc(func(e *adventuria.OnBeforeRollEvent) (*event.Result, error) {
			dicesAny, _ := ef.DecodeValue(ef.GetString("value"))
			dices := dicesAny.([]string)
			diceList := make([]*adventuria.Dice, len(dices))
			for i, name := range dices {
				diceList[i] = adventuria.DiceList[name]
			}
			e.Dices = diceList

			callback(e.AppContext)

			return e.Next()
		}),
	}, nil
}

func (ef *ChangeDicesEffect) Verify(_ adventuria.AppContext, value string) error {
	names, err := ef.DecodeValue(value)
	if err != nil {
		return fmt.Errorf("changeDices: %w", err)
	}

	for _, name := range names.([]string) {
		if _, ok := adventuria.DiceList[name]; !ok {
			return errors.New("changeDices: unknown dice")
		}
	}

	return nil
}

func (ef *ChangeDicesEffect) DecodeValue(value string) (any, error) {
	return strings.Split(value, ";"), nil
}

func (ef *ChangeDicesEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

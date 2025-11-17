package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"errors"
	"fmt"
	"strings"
)

type ChangeDicesEffect struct {
	adventuria.EffectBase
}

func (ef *ChangeDicesEffect) Subscribe(
	user adventuria.User,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		user.OnBeforeRoll().BindFunc(func(e *adventuria.OnBeforeRollEvent) error {
			dicesAny, _ := ef.DecodeValue(ef.GetString("value"))
			dices := dicesAny.([]string)
			diceList := make([]*adventuria.Dice, len(dices))
			for i, name := range dices {
				diceList[i] = adventuria.DiceList[name]
			}
			e.Dices = diceList

			callback()

			return e.Next()
		}),
	}
}

func (ef *ChangeDicesEffect) Verify(value string) error {
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
	return strings.Split(value, ","), nil
}

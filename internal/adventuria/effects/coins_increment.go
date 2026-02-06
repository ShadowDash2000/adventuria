package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type CoinsIncrementEffect struct {
	adventuria.EffectRecord
}

func (ef *CoinsIncrementEffect) CanUse(_ adventuria.AppContext, _ adventuria.EffectContext) bool {
	return true
}

func (ef *CoinsIncrementEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	value, err := ef.DecodeValue(ef.GetString("value"))
	if err != nil {
		return nil, err
	}

	switch value.Event {
	case "onAfterItemSave":
		return []event.Unsubscribe{
			ctx.User.OnAfterItemSave().BindFunc(func(e *adventuria.OnAfterItemSave) (*event.Result, error) {
				if e.Item.IDInventory() == ctx.InvItemID {
					err = ctx.User.AddBalance(e.AppContext, value.Value)
					if err != nil {
						return &event.Result{
							Success: false,
							Error:   "internal error: can't update user balance",
						}, fmt.Errorf("coinsIncrementEffect: %w", err)
					}
					callback(e.AppContext)
				}

				return e.Next()
			}),
		}, nil
	default:
		return nil, nil
	}
}

func (ef *CoinsIncrementEffect) Verify(_ adventuria.AppContext, value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

type CoinsIncrementEffectValue struct {
	Event string
	Value int
}

func (ef *CoinsIncrementEffect) DecodeValue(value string) (*CoinsIncrementEffectValue, error) {
	values := strings.Split(value, ";")
	if len(values) != 2 {
		return nil, fmt.Errorf("coinsIncrement: invalid value, expected format 'value;event': %s", value)
	}

	coins, err := strconv.Atoi(values[0])
	if err != nil {
		return nil, fmt.Errorf("coinsIncrement: invalid value: %s", values[1])
	}

	events := []string{
		"onAfterItemSave",
	}

	if !slices.Contains(events, values[1]) {
		return nil, fmt.Errorf("coinsIncrement: invalid event: %s", values[0])
	}

	return &CoinsIncrementEffectValue{
		Event: values[1],
		Value: coins,
	}, nil
}

func (ef *CoinsIncrementEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

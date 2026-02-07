package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"fmt"
	"strconv"
	"strings"
)

type CoinsForAllEffect struct {
	adventuria.EffectRecord
}

func (ef *CoinsForAllEffect) CanUse(_ adventuria.AppContext, _ adventuria.EffectContext) bool {
	return true
}

func (ef *CoinsForAllEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	decodedValue, err := ef.DecodeValue(ef.GetString("value"))
	if err != nil {
		return nil, err
	}

	return []event.Unsubscribe{
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*event.Result, error) {
			if e.InvItemId == ctx.InvItemID {
				ctx.User.AddBalance(decodedValue.CoinsForPlayer)

				for _, user := range adventuria.GameUsers.GetAll() {
					if user.ID() == ctx.User.ID() {
						continue
					}
					user.Lock()
					user.AddBalance(decodedValue.CoinsForOther)
					user.Unlock()
				}

				callback(e.AppContext)
			}

			return e.Next()
		}),
	}, nil
}

func (ef *CoinsForAllEffect) Verify(_ adventuria.AppContext, value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

type CoinsForAllEffectValue struct {
	CoinsForPlayer int
	CoinsForOther  int
}

func (ef *CoinsForAllEffect) DecodeValue(value string) (*CoinsForAllEffectValue, error) {
	values := strings.Split(value, ";")
	if len(values) != 2 {
		return nil, fmt.Errorf("coinsForAll: invalid value, expected format 'event;value': %s", value)
	}

	var err error
	coins := make([]int, len(values))
	for i, value := range values {
		coins[i], err = strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("coinsIncrement: invalid value: %s", value)
		}
	}

	return &CoinsForAllEffectValue{
		CoinsForPlayer: coins[0],
		CoinsForOther:  coins[1],
	}, nil
}

func (ef *CoinsForAllEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

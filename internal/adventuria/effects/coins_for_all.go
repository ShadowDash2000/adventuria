package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/event"
	"fmt"
	"strconv"
	"strings"

	"github.com/pocketbase/dbx"
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
				ctx.User.SetBalance(ctx.User.Balance() + decodedValue.CoinsForPlayer)

				query := fmt.Sprintf(
					"UPDATE %s SET %[2]s = %[2]s + {:coins} WHERE id != {:currentUserId}",
					schema.CollectionUsers,
					schema.UserSchema.Balance,
				)
				_, err = e.App.DB().
					NewQuery(query).
					Bind(dbx.Params{
						"coins":         decodedValue.CoinsForOther,
						"currentUserId": ctx.User.ID(),
					}).
					Execute()
				if err != nil {
					return &event.Result{
						Success: false,
						Error:   "internal error: can't update user balance",
					}, fmt.Errorf("coinsForAllEffect: %w", err)
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

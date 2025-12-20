package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"adventuria/pkg/event"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type ChangeMaxGamePriceEffect struct {
	adventuria.EffectBase
}

func (ef *ChangeMaxGamePriceEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	valAny, err := ef.DecodeValue(ef.GetString("value"))
	if err != nil {
		return nil, err
	}

	val := valAny.(ChangeMaxGamePriceEffectValue)

	switch val.Type {
	case "usable":
		return []event.Unsubscribe{
			ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) error {
				if e.InvItemId == ctx.InvItemID {
					ok, err := ef.tryToApplyEffect(ctx.User)
					if err != nil {
						return fmt.Errorf("changeMaxGamePrice: %w", err)
					}

					if ok {
						callback()
					}
				}

				return e.Next()
			}),
		}, nil
	case "unusable":
		return []event.Unsubscribe{
			ctx.User.OnAfterMove().BindFunc(func(e *adventuria.OnAfterMoveEvent) error {
				ok, err := ef.tryToApplyEffect(ctx.User)
				if err != nil {
					return fmt.Errorf("changeMaxGamePrice: %w", err)
				}

				if ok {
					callback()
				}

				return e.Next()
			}),
			ctx.User.OnAfterItemSave().BindFunc(func(e *adventuria.OnAfterItemSave) error {
				if e.Item.IDInventory() != ctx.InvItemID {
					return e.Next()
				}

				ok, err := ef.tryToApplyEffect(ctx.User)
				if err != nil {
					return fmt.Errorf("changeMaxGamePrice: %w", err)
				}

				if ok {
					callback()
				}

				return e.Next()
			}),
		}, nil
	default:
		return nil, nil
	}
}

func (ef *ChangeMaxGamePriceEffect) tryToApplyEffect(user adventuria.User) (bool, error) {
	if !adventuria.GameActions.CanDo(user, "rollWheel") {
		return false, nil
	}

	cell, ok := user.CurrentCell()
	if !ok {
		return false, nil
	}

	cellGame, ok := cell.(*cells.CellGame)
	if !ok {
		return false, nil
	}

	valAny, err := ef.DecodeValue(ef.GetString("value"))
	if err != nil {
		return false, err
	}

	val := valAny.(ChangeMaxGamePriceEffectValue)

	user.LastAction().CustomGameFilter().MaxPrice = val.Price
	if err = cellGame.CheckCustomFilter(user); err != nil {
		return false, err
	}

	return true, nil
}

func (ef *ChangeMaxGamePriceEffect) Verify(value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

type ChangeMaxGamePriceEffectValue struct {
	Price int
	Type  string
}

func (ef *ChangeMaxGamePriceEffect) DecodeValue(value string) (any, error) {
	vals := strings.Split(value, ";")
	if len(vals) != 2 {
		return nil, fmt.Errorf("changeMaxGamePrice: invalid value: %s", value)
	}

	var (
		res   ChangeMaxGamePriceEffectValue
		err   error
		types = []string{"usable", "unusable"}
	)

	res.Price, err = strconv.Atoi(vals[0])
	if err != nil {
		return nil, fmt.Errorf("changeMaxGamePrice: invalid value: %s", value)
	}

	if !slices.Contains(types, vals[1]) {
		return nil, fmt.Errorf("changeMaxGamePrice: invalid event: %s", vals[1])
	}
	res.Type = vals[1]

	return res, nil
}

package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/event"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type ChangeMaxGamePriceEffect struct {
	adventuria.EffectRecord
}

func (ef *ChangeMaxGamePriceEffect) CanUse(appCtx adventuria.AppContext, ctx adventuria.EffectContext) bool {
	if !adventuria.GameActions.CanDo(appCtx, ctx.User, "rollWheel") {
		return false
	}

	cell, ok := ctx.User.CurrentCell()
	if !ok {
		return false
	}

	if cell.Type() != "game" {
		return false
	}

	if cell.IsCustomFilterNotAllowed() {
		return false
	}

	if filterId := cell.Filter(); filterId != "" {
		filterRecord, err := appCtx.App.FindRecordById(
			schema.CollectionActivityFilter,
			filterId,
		)
		if err != nil {
			return false
		}

		if len(filterRecord.GetStringSlice("activities")) > 0 {
			return false
		}
	}

	return true
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
			ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*event.Result, error) {
				if e.InvItemId == ctx.InvItemID {
					res, err := ef.tryToApplyEffect(e.AppContext, ctx.User)
					if err != nil {
						return res, err
					}

					if res.Success {
						callback(e.AppContext)
					} else {
						return res, nil
					}
				}

				return e.Next()
			}),
		}, nil
	case "unusable":
		return []event.Unsubscribe{
			ctx.User.OnAfterMove().BindFunc(func(e *adventuria.OnAfterMoveEvent) (*event.Result, error) {
				if !ef.CanUse(e.AppContext, ctx) {
					return e.Next()
				}

				res, err := ef.tryToApplyEffect(e.AppContext, ctx.User)
				if err != nil {
					return res, err
				}

				if res.Success {
					callback(e.AppContext)
				} else {
					return res, nil
				}

				return e.Next()
			}),
			ctx.User.OnAfterItemSave().BindFunc(func(e *adventuria.OnAfterItemSave) (*event.Result, error) {
				if e.Item.IDInventory() != ctx.InvItemID {
					return e.Next()
				}

				if !ef.CanUse(e.AppContext, ctx) {
					return e.Next()
				}

				res, err := ef.tryToApplyEffect(e.AppContext, ctx.User)
				if err != nil {
					return res, err
				}

				if res.Success {
					callback(e.AppContext)
				} else {
					return res, nil
				}

				return e.Next()
			}),
		}, nil
	default:
		return nil, nil
	}
}

func (ef *ChangeMaxGamePriceEffect) tryToApplyEffect(ctx adventuria.AppContext, user adventuria.User) (*event.Result, error) {
	cell, ok := user.CurrentCell()
	if !ok {
		return &event.Result{
			Success: false,
			Error:   "current cell not found",
		}, nil
	}

	cellGame, ok := cell.(adventuria.CellWheel)
	if !ok {
		return &event.Result{
			Success: false,
			Error:   "current cell isn't game cell",
		}, nil
	}

	valAny, err := ef.DecodeValue(ef.GetString("value"))
	if err != nil {
		return &event.Result{
			Success: false,
			Error:   "internal error: invalid value in \"change_max_game_price\" effect",
		}, fmt.Errorf("changeMaxGamePrice: %w", err)
	}

	val := valAny.(ChangeMaxGamePriceEffectValue)

	filter, err := user.LastAction().CustomActivityFilter()
	if err != nil {
		return &event.Result{
			Success: false,
			Error:   "internal error: can't get custom activity filter",
		}, fmt.Errorf("changeMaxGamePrice: %w", err)
	}

	filter.MaxPrice = val.Price
	filter.MinPrice = -1
	user.LastAction().SetCustomActivityFilter(*filter)

	if err = cellGame.RefreshItems(ctx, user); err != nil {
		return &event.Result{
			Success: false,
			Error:   "internal error: can't refresh cell items in \"change_max_game_price\" effect",
		}, fmt.Errorf("changeMaxGamePrice: %w", err)
	}

	return &event.Result{
		Success: true,
	}, nil
}

func (ef *ChangeMaxGamePriceEffect) Verify(_ adventuria.AppContext, value string) error {
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

func (ef *ChangeMaxGamePriceEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

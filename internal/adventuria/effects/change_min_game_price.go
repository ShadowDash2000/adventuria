package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/event"
	"fmt"
	"strconv"
)

type ChangeMinGamePriceEffect struct {
	adventuria.EffectRecord
}

func (ef *ChangeMinGamePriceEffect) CanUse(appCtx adventuria.AppContext, ctx adventuria.EffectContext) bool {
	if ok := adventuria.GameActions.CanDo(appCtx, ctx.User, "rollWheel"); !ok {
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

func (ef *ChangeMinGamePriceEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*event.Result, error) {
			if e.InvItemId == ctx.InvItemID {
				cell, ok := ctx.User.CurrentCell()
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

				if i := ef.GetInt("value"); i != 0 {
					ctx.User.LastAction().CustomActivityFilter().MinPrice = i
					ctx.User.LastAction().CustomActivityFilter().MaxPrice = -1
					if err := cellGame.RefreshItems(e.AppContext, ctx.User); err != nil {
						return &event.Result{
							Success: false,
							Error:   "internal error: can't refresh cell items in \"change_min_game_price\" effect",
						}, fmt.Errorf("changeMinGamePrice: %w", err)
					}

					callback(e.AppContext)
				}
			}

			return e.Next()
		}),
	}, nil
}

func (ef *ChangeMinGamePriceEffect) Verify(_ adventuria.AppContext, value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *ChangeMinGamePriceEffect) DecodeValue(value string) (any, error) {
	return strconv.Atoi(value)
}

func (ef *ChangeMinGamePriceEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

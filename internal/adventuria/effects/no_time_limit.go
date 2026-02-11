package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
	"errors"
	"fmt"
)

type NoTimeLimitEffect struct {
	adventuria.EffectRecord
}

func (ef *NoTimeLimitEffect) CanUse(appCtx adventuria.AppContext, ctx adventuria.EffectContext) bool {
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

	return true
}

func (ef *NoTimeLimitEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterMove().BindFunc(func(e *adventuria.OnAfterMoveEvent) (*result.Result, error) {
			if !ef.CanUse(e.AppContext, ctx) {
				return e.Next()
			}

			res, err := ef.tryToApplyEffect(e.AppContext, ctx.User)
			if err != nil {
				return res, err
			}

			if res.Ok() {
				callback(e.AppContext)
			} else {
				return res, nil
			}

			return e.Next()
		}),
		ctx.User.OnAfterItemSave().BindFunc(func(e *adventuria.OnAfterItemSave) (*result.Result, error) {
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

			if res.Ok() {
				callback(e.AppContext)
			} else {
				return res, nil
			}

			return e.Next()
		}),
	}, nil
}

func (ef *NoTimeLimitEffect) tryToApplyEffect(ctx adventuria.AppContext, user adventuria.User) (*result.Result, error) {
	cell, ok := user.CurrentCell()
	if !ok {
		return result.Err("internal error: current cell not found"),
			errors.New("noTimeLimit: current cell not found")
	}

	cellGame, ok := cell.(adventuria.CellWheel)
	if !ok {
		return result.Err("current cell isn't wheel cell"), nil
	}

	filter, err := user.LastAction().CustomActivityFilter()
	if err != nil {
		return result.Err("internal error: can't get custom activity filter"),
			fmt.Errorf("noTimeLimit: %w", err)
	}

	filter.MinCampaignTime = -1
	filter.MaxCampaignTime = -1
	if err := cellGame.RefreshItems(ctx, user); err != nil {
		return result.Err("internal error: can't refresh cell items"),
			fmt.Errorf("noTimeLimit: %w", err)
	}

	user.LastAction().SetCustomActivityFilter(*filter)

	return result.Ok(), nil
}

func (ef *NoTimeLimitEffect) Verify(_ adventuria.AppContext, _ string) error {
	return nil
}

func (ef *NoTimeLimitEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

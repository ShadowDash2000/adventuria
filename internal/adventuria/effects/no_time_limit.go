package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"adventuria/pkg/event"
	"fmt"
)

type NoTimeLimitEffect struct {
	adventuria.EffectRecord
}

func (ef *NoTimeLimitEffect) CanUse(_ adventuria.EffectContext) bool {
	return true
}

func (ef *NoTimeLimitEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterMove().BindFunc(func(e *adventuria.OnAfterMoveEvent) (*event.Result, error) {
			res, err := ef.tryToApplyEffect(ctx.User)
			if err != nil {
				return res, err
			}

			if res.Success {
				callback()
			} else {
				return res, nil
			}

			return e.Next()
		}),
		ctx.User.OnAfterItemSave().BindFunc(func(e *adventuria.OnAfterItemSave) (*event.Result, error) {
			if e.Item.IDInventory() != ctx.InvItemID {
				return e.Next()
			}

			res, err := ef.tryToApplyEffect(ctx.User)
			if err != nil {
				return res, err
			}

			if res.Success {
				callback()
			} else {
				return res, nil
			}

			return e.Next()
		}),
	}, nil
}

func (ef *NoTimeLimitEffect) tryToApplyEffect(user adventuria.User) (*event.Result, error) {
	if !adventuria.GameActions.CanDo(user, "rollWheel") {
		return &event.Result{
			Success: false,
			Error:   "user can't perform roll wheel action",
		}, nil
	}

	cell, ok := user.CurrentCell()
	if !ok {
		return &event.Result{
			Success: false,
			Error:   "current cell not found",
		}, nil
	}

	cellGame, ok := cell.(*cells.CellGame)
	if !ok {
		return &event.Result{
			Success: false,
			Error:   "current cell isn't game cell",
		}, nil
	}

	if cell.Type() != "game" {
		return &event.Result{
			Success: false,
			Error:   "current cell isn't game cell",
		}, nil
	}

	user.LastAction().CustomActivityFilter().MinCampaignTime = -1
	user.LastAction().CustomActivityFilter().MaxCampaignTime = -1
	if err := cellGame.RefreshItems(user); err != nil {
		return &event.Result{
			Success: false,
			Error:   "internal error: can't refresh cell items in \"no_time_limit\" effect",
		}, fmt.Errorf("noTimeLimit: %w", err)
	}

	return &event.Result{
		Success: true,
	}, nil
}

func (ef *NoTimeLimitEffect) Verify(_ string) error {
	return nil
}

func (ef *NoTimeLimitEffect) GetVariants(_ adventuria.EffectContext) any {
	return nil
}

package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"adventuria/pkg/event"
	"fmt"
)

type NoTimeLimitEffect struct {
	adventuria.EffectBase
}

func (ef *NoTimeLimitEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterMove().BindFunc(func(e *adventuria.OnAfterMoveEvent) error {
			ok, err := ef.tryToApplyEffect(ctx.User)
			if err != nil {
				return fmt.Errorf("noTimeLimit: %w", err)
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
				return fmt.Errorf("noTimeLimit: %w", err)
			}

			if ok {
				callback()
			}

			return e.Next()
		}),
	}, nil
}

func (ef *NoTimeLimitEffect) tryToApplyEffect(user adventuria.User) (bool, error) {
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

	user.LastAction().CustomGameFilter().MinCampaignTime = -1
	user.LastAction().CustomGameFilter().MaxCampaignTime = -1
	if err := cellGame.CheckCustomFilter(user); err != nil {
		return false, err
	}

	return true, nil
}

func (ef *NoTimeLimitEffect) Verify(_ string) error {
	return nil
}

func (ef *NoTimeLimitEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}

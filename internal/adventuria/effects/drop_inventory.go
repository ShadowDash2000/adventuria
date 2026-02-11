package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
	"fmt"
)

type DropInventoryEffect struct {
	adventuria.EffectRecord
}

func (ef *DropInventoryEffect) CanUse(_ adventuria.AppContext, _ adventuria.EffectContext) bool {
	return true
}

func (ef *DropInventoryEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemSave().BindFunc(func(e *adventuria.OnAfterItemSave) (*result.Result, error) {
			if e.Item.IDInventory() == ctx.InvItemID {
				err := ctx.User.Inventory().DropInventory(e.AppContext)
				if err != nil {
					return result.Err("internal error: failed to drop inventory"),
						fmt.Errorf("dropInventory: %w", err)
				}

				callback(e.AppContext)
			}

			return e.Next()
		}),
	}, nil
}

func (ef *DropInventoryEffect) Verify(_ adventuria.AppContext, _ string) error {
	return nil
}

func (ef *DropInventoryEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

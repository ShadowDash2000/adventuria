package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"fmt"
)

type DropInventoryEffect struct {
	adventuria.EffectBase
}

func (ef *DropInventoryEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemSave().BindFunc(func(e *adventuria.OnAfterItemSave) (*event.Result, error) {
			if e.Item.IDInventory() == ctx.InvItemID {
				err := ctx.User.Inventory().DropInventory()
				if err != nil {
					return &event.Result{
						Success: false,
						Error:   "internal error: can't drop inventory",
					}, fmt.Errorf("dropInventory: %w", err)
				}

				callback()
			}

			return e.Next()
		}),
	}, nil
}

func (ef *DropInventoryEffect) Verify(_ string) error {
	return nil
}

func (ef *DropInventoryEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}

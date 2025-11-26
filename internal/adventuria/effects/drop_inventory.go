package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type DropInventoryEffect struct {
	adventuria.EffectBase
}

func (ef *DropInventoryEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		ctx.User.OnAfterAction().BindFunc(func(e *adventuria.OnAfterActionEvent) error {
			err := ctx.User.Inventory().DropInventory()
			if err != nil {
				return err
			}

			callback()

			return e.Next()
		}),
	}
}

func (ef *DropInventoryEffect) Verify(_ string) error {
	return nil
}

func (ef *DropInventoryEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}

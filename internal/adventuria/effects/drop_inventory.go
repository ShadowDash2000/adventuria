package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type DropInventoryEffect struct {
	adventuria.EffectBase
}

func (ef *DropInventoryEffect) Subscribe(
	user adventuria.User,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		user.OnAfterAction().BindFunc(func(e *adventuria.OnAfterActionEvent) error {
			err := user.Inventory().DropInventory()
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

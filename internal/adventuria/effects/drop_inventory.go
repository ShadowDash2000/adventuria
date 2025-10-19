package effects

import (
	"adventuria/internal/adventuria"
)

type DropInventoryEffect struct {
	adventuria.EffectBase
}

func (ef *DropInventoryEffect) Subscribe(callback adventuria.EffectCallback) {
	ef.PoolUnsubscribers(
		ef.User().OnAfterAction().BindFunc(func(e *adventuria.OnAfterActionEvent) error {
			err := ef.User().Inventory().DropInventory()
			if err != nil {
				return err
			}

			callback()

			return e.Next()
		}),
	)
}

func (ef *DropInventoryEffect) Verify(_ string) error {
	return nil
}

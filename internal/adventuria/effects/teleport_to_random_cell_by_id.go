package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/helper"
)

type TeleportToRandomCellByIdEffect struct {
	adventuria.EffectBase
}

func (ef *TeleportToRandomCellByIdEffect) Subscribe(callback adventuria.EffectCallback) {
	ef.PoolUnsubscribers(
		ef.User().OnAfterAction().BindFunc(func(e *adventuria.OnAfterActionEvent) error {
			cellIds := []string{} // TODO: get cell ids from effect record
			if len(cellIds) > 0 {
				cellId := helper.RandomItemFromSlice(cellIds)
				err := ef.User().MoveToCellId(cellId)
				if err != nil {
					return err
				}
			}

			callback()

			return e.Next()
		}),
	)
}

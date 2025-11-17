package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/helper"
)

type TeleportToRandomCellByIdEffect struct {
	adventuria.EffectBase
}

func (ef *TeleportToRandomCellByIdEffect) Subscribe(
	user adventuria.User,
	callback adventuria.EffectCallback,
) []event.Unsubscribe {
	return []event.Unsubscribe{
		user.OnAfterAction().BindFunc(func(e *adventuria.OnAfterActionEvent) error {
			cellIds := []string{} // TODO: get cell ids from effect record
			if len(cellIds) > 0 {
				cellId := helper.RandomItemFromSlice(cellIds)
				err := user.MoveToCellId(cellId)
				if err != nil {
					return err
				}
			}

			callback()

			return e.Next()
		}),
	}
}

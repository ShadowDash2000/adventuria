package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type GoToJailEffect struct {
	adventuria.EffectBase
}

func (ef *GoToJailEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemSave().BindFunc(func(e *adventuria.OnAfterItemSave) error {
			if e.Item.IDInventory() == ctx.InvItemID {
				_, err := ctx.User.MoveToClosestCellType("jail")
				if err != nil {
					return err
				}

				ctx.User.SetIsInJail(true)

				err = ctx.User.OnAfterGoToJail().Trigger(&adventuria.OnAfterGoToJailEvent{})
				if err != nil {
					adventuria.PocketBase.Logger().Error(
						"goToJailEffect: failed to trigger onAfterGoToJail event",
						"error",
						err,
					)
				}

				callback()
			}

			return e.Next()
		}),
	}, nil
}

func (ef *GoToJailEffect) Verify(_ string) error {
	return nil
}

func (ef *GoToJailEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}

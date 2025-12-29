package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"fmt"
)

type GoToJailEffect struct {
	adventuria.EffectRecord
}

func (ef *GoToJailEffect) CanUse(_ adventuria.EffectContext) bool {
	return true
}

func (ef *GoToJailEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemSave().BindFunc(func(e *adventuria.OnAfterItemSave) (*event.Result, error) {
			if e.Item.IDInventory() == ctx.InvItemID {
				_, err := ctx.User.MoveToClosestCellType("jail")
				if err != nil {
					return &event.Result{
						Success: false,
						Error:   "internal error: can't move to jail cell",
					}, fmt.Errorf("goToJailEffect: %w", err)
				}

				ctx.User.SetIsInJail(true)

				callback()

				res, err := ctx.User.OnAfterGoToJail().Trigger(&adventuria.OnAfterGoToJailEvent{})
				if err != nil {
					return res, err
				}
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

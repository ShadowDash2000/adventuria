package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type ReturnToPrevCellEffect struct {
	adventuria.EffectBase
}

func (ef *ReturnToPrevCellEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) error {
			if e.InvItemId != ctx.InvItemID {
				return e.Next()
			}

			latestDiceRoll := ctx.User.LastAction().DiceRoll()
			if latestDiceRoll == 0 {
				return e.Next()
			}

			_, err := ctx.User.Move(-latestDiceRoll)
			if err != nil {
				return err
			}

			ctx.User.LastAction().SetCanMove(true)

			callback()

			return e.Next()
		}),
	}, nil
}

func (ef *ReturnToPrevCellEffect) Verify(_ string) error {
	return nil
}

func (ef *ReturnToPrevCellEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}

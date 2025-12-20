package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"errors"
	"fmt"
	"slices"
)

type ChooseGameEffect struct {
	adventuria.EffectBase
}

func (ef *ChooseGameEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) error {
			if ctx.InvItemID == e.InvItemId {
				if ok := adventuria.GameActions.CanDo(ctx.User, "done"); !ok {
					return errors.New("chooseGame: user can't do done action")
				}

				if gameId, ok := e.Request["game_id"].(string); ok {
					itemsList, err := ctx.User.LastAction().ItemsList()
					if err != nil {
						return fmt.Errorf("chooseGame: %w", err)
					}

					if !slices.Contains(itemsList, gameId) {
						return errors.New("chooseGame: game_id not found in items list")
					}

					ctx.User.LastAction().SetGame(gameId)

					callback()
				} else {
					return errors.New("chooseGame: game_id not found")
				}
			}

			return e.Next()
		}),
	}, nil
}

func (ef *ChooseGameEffect) Verify(_ string) error {
	return nil
}

func (ef *ChooseGameEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}

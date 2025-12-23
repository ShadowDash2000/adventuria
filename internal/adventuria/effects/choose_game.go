package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
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
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*event.Result, error) {
			if ctx.InvItemID == e.InvItemId {
				if ok := adventuria.GameActions.CanDo(ctx.User, "done"); !ok {
					return &event.Result{
						Success: false,
						Error:   "user can't perform done action",
					}, nil
				}

				if gameId, ok := e.Request["game_id"].(string); ok {
					itemsList, err := ctx.User.LastAction().ItemsList()
					if err != nil {
						return &event.Result{
							Success: false,
							Error:   "internal error: can't unmarshal items list in \"choose_game\" effect",
						}, fmt.Errorf("chooseGame: %w", err)
					}

					if !slices.Contains(itemsList, gameId) {
						return &event.Result{
							Success: false,
							Error:   "game_id not found in items list",
						}, nil
					}

					ctx.User.LastAction().SetGame(gameId)

					callback()
				} else {
					return &event.Result{
						Success: false,
						Error:   "request error: game_id not specified",
					}, nil
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

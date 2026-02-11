package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
	"fmt"
	"slices"
)

type ChooseGameEffect struct {
	adventuria.EffectRecord
}

func (ef *ChooseGameEffect) CanUse(appCtx adventuria.AppContext, ctx adventuria.EffectContext) bool {
	return adventuria.GameActions.CanDo(appCtx, ctx.User, "done")
}

func (ef *ChooseGameEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*result.Result, error) {
			if ctx.InvItemID == e.InvItemId {
				if gameId, ok := e.Data["game_id"].(string); ok {
					itemsList, err := ctx.User.LastAction().ItemsList()
					if err != nil {
						return result.Err("internal error: failed to decode action's items list"),
							fmt.Errorf("chooseGame: %w", err)
					}

					if !slices.Contains(itemsList, gameId) {
						return result.Err("game_id not found in items list"), nil
					}

					ctx.User.LastAction().SetActivity(gameId)

					callback(e.AppContext)
				} else {
					return result.Err("game_id not specified"), nil
				}
			}

			return e.Next()
		}),
	}, nil
}

func (ef *ChooseGameEffect) Verify(_ adventuria.AppContext, _ string) error {
	return nil
}

func (ef *ChooseGameEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

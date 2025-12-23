package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"errors"
)

type ChangeGameByIdEffect struct {
	adventuria.EffectBase
}

func (ef *ChangeGameByIdEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	if ok := adventuria.GameActions.CanDo(ctx.User, "done"); !ok {
		return nil, errors.New("changeGameById: user can't do done action")
	}

	return []event.Unsubscribe{
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*event.Result, error) {
			if e.InvItemId == ctx.InvItemID {
				if ok := adventuria.GameActions.CanDo(ctx.User, "done"); !ok {
					return &event.Result{
						Success: false,
						Error:   "user can't perform done action",
					}, nil
				}

				ctx.User.LastAction().SetGame(ef.GetString("value"))

				callback()
			}

			return e.Next()
		}),
		ctx.User.OnAfterMove().BindFunc(func(e *adventuria.OnAfterMoveEvent) (*event.Result, error) {
			callback()
			return e.Next()
		}),
	}, nil
}

func (ef *ChangeGameByIdEffect) Verify(gameId string) error {
	_, err := adventuria.PocketBase.FindRecordById(
		adventuria.GameCollections.Get(adventuria.CollectionGames),
		gameId,
	)
	return err
}

func (ef *ChangeGameByIdEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}

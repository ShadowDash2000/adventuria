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
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) error {
			if e.InvItemId == ctx.InvItemID {
				if ok := adventuria.GameActions.CanDo(ctx.User, "done"); !ok {
					return errors.New("changeGameById: user can't do done action")
				}

				ctx.User.LastAction().SetGame(ef.GetString("value"))

				callback()
			}

			return e.Next()
		}),
		ctx.User.OnAfterMove().BindFunc(func(e *adventuria.OnAfterMoveEvent) error {
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

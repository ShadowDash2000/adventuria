package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
)

type ChangeGameByIdEffect struct {
	adventuria.EffectRecord
}

func (ef *ChangeGameByIdEffect) CanUse(ctx adventuria.EffectContext) bool {
	if ok := adventuria.GameActions.CanDo(ctx.User, "drop"); !ok {
		return false
	}

	if ok := adventuria.GameActions.CanDo(ctx.User, "done"); !ok {
		return false
	}

	return true
}

func (ef *ChangeGameByIdEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*event.Result, error) {
			if e.InvItemId == ctx.InvItemID {
				ctx.User.LastAction().SetActivity(ef.GetString("value"))
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
		adventuria.GameCollections.Get(adventuria.CollectionActivities),
		gameId,
	)
	return err
}

func (ef *ChangeGameByIdEffect) GetVariants(_ adventuria.EffectContext) any {
	return nil
}

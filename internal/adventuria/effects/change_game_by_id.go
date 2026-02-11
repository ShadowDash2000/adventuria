package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
)

type ChangeGameByIdEffect struct {
	adventuria.EffectRecord
}

func (ef *ChangeGameByIdEffect) CanUse(appCtx adventuria.AppContext, ctx adventuria.EffectContext) bool {
	if ok := adventuria.GameActions.CanDo(appCtx, ctx.User, "drop"); !ok {
		return false
	}

	if ok := adventuria.GameActions.CanDo(appCtx, ctx.User, "done"); !ok {
		return false
	}

	return true
}

func (ef *ChangeGameByIdEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*result.Result, error) {
			if e.InvItemId == ctx.InvItemID {
				ctx.User.LastAction().SetActivity(ef.GetString("value"))
				callback(e.AppContext)
			}

			return e.Next()
		}),
		ctx.User.OnAfterMove().BindFunc(func(e *adventuria.OnAfterMoveEvent) (*result.Result, error) {
			callback(e.AppContext)
			return e.Next()
		}),
	}, nil
}

func (ef *ChangeGameByIdEffect) Verify(ctx adventuria.AppContext, gameId string) error {
	_, err := ctx.App.FindRecordById(
		adventuria.GameCollections.Get(schema.CollectionActivities),
		gameId,
	)
	return err
}

func (ef *ChangeGameByIdEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

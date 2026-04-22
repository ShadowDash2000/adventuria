package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
)

var _ adventuria.EffectVerifiable = (*ChangeGameByIdEffect)(nil)

type ChangeGameByIdEffect struct {
	adventuria.EffectRecord
}

func (ef *ChangeGameByIdEffect) CanUse(appCtx adventuria.AppContext, ctx adventuria.EffectContext) bool {
	cell, ok := ctx.Player.Progress().CurrentCell()
	if !ok {
		return false
	}

	if ok = cell.IsChangeGameNotAllowed(); ok {
		return false
	}

	if ok = adventuria.GameActions.CanDo(appCtx, ctx.Player, "drop"); !ok {
		return false
	}

	if ok = adventuria.GameActions.CanDo(appCtx, ctx.Player, "done"); !ok {
		return false
	}

	return true
}

func (ef *ChangeGameByIdEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.Player.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*result.Result, error) {
			if e.InvItemId == ctx.InvItemID {
				ctx.Player.LastAction().SetActivity(ef.GetString("value"))
				callback(e.AppContext)
			}

			return e.Next()
		}),
		ctx.Player.OnAfterMove().BindFunc(func(e *adventuria.OnAfterMoveEvent) (*result.Result, error) {
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

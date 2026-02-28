package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
	"errors"
	"slices"
)

type NothingEffect struct {
	adventuria.EffectRecord
}

func (ef *NothingEffect) CanUse(_ adventuria.AppContext, _ adventuria.EffectContext) bool {
	return true
}

func (ef *NothingEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	switch ef.GetString("value") {
	case "onAfterItemAdd":
		return []event.Unsubscribe{
			ctx.User.OnAfterItemAdd().BindFunc(func(e *adventuria.OnAfterItemAdd) (*result.Result, error) {
				callback(e.AppContext)
				return e.Next()
			}),
		}, nil
	case "onAfterItemUse":
		return []event.Unsubscribe{
			ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*result.Result, error) {
				if e.InvItemId == ctx.InvItemID {
					callback(e.AppContext)
				}
				return e.Next()
			}),
		}, nil
	case "onBeforeGameDone":
		return []event.Unsubscribe{
			ctx.User.OnBeforeDone().BindFunc(func(e *adventuria.OnBeforeDoneEvent) (*result.Result, error) {
				currentCell, ok := ctx.User.CurrentCell()
				if !ok {
					return result.Err("internal error: current cell not found"),
						errors.New("nothing: current cell not found")
				}

				if currentCell.InCategories([]string{"activity", "game"}) {
					callback(e.AppContext)
				}

				return e.Next()
			}),
		}, nil
	default:
		return nil, nil
	}
}

func (ef *NothingEffect) Verify(_ adventuria.AppContext, v string) error {
	events := []string{
		"onAfterItemAdd",
		"onAfterItemUse",
		"onBeforeGameDone",
	}

	if ok := slices.Contains(events, v); !ok {
		return errors.New("unknown event")
	}

	return nil
}

func (ef *NothingEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

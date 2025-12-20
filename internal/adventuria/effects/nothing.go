package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"errors"
	"slices"
)

type NothingEffect struct {
	adventuria.EffectBase
}

func (ef *NothingEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	switch ef.GetString("value") {
	case "onAfterItemAdd":
		return []event.Unsubscribe{
			ctx.User.OnAfterItemAdd().BindFunc(func(e *adventuria.OnAfterItemAdd) error {
				callback()
				return e.Next()
			}),
		}, nil
	case "onAfterItemUse":
		return []event.Unsubscribe{
			ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) error {
				if e.InvItemId == ctx.InvItemID {
					callback()
				}
				return e.Next()
			}),
		}, nil
	case "onBeforeDone":
		return []event.Unsubscribe{
			ctx.User.OnBeforeDone().BindFunc(func(e *adventuria.OnBeforeDoneEvent) error {
				callback()
				return e.Next()
			}),
		}, nil
	default:
		return nil, nil
	}
}

func (ef *NothingEffect) Verify(v string) error {
	events := []string{
		"onAfterItemAdd",
		"onAfterItemUse",
		"onBeforeDone",
	}

	if ok := slices.Contains(events, v); !ok {
		return errors.New("unknown event")
	}

	return nil
}

func (ef *NothingEffect) DecodeValue(_ string) (any, error) {
	return nil, nil
}

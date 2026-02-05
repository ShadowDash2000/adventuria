package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"fmt"
	"slices"
)

type GoToJailEffect struct {
	adventuria.EffectRecord
}

func (ef *GoToJailEffect) CanUse(appCtx adventuria.AppContext, ctx adventuria.EffectContext) bool {
	canRollWheel := adventuria.GameActions.CanDo(appCtx, ctx.User, "rollWheel")
	if canRollWheel {
		return false
	}

	canDone := adventuria.GameActions.CanDo(appCtx, ctx.User, "done")
	canDrop := adventuria.GameActions.CanDo(appCtx, ctx.User, "drop")

	if canDone && !canDrop {
		return false
	}

	return true
}

func (ef *GoToJailEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	decodedValue, err := ef.DecodeValue(ef.GetString("value"))
	if err != nil {
		return nil, err
	}

	switch decodedValue.Event {
	case "onAfterItemSave":
		return []event.Unsubscribe{
			ctx.User.OnAfterItemSave().BindFunc(func(e *adventuria.OnAfterItemSave) (*event.Result, error) {
				if e.Item.IDInventory() != ctx.InvItemID {
					return e.Next()
				}

				res, err := ef.tryToApplyEffect(e.AppContext, ctx.User)
				if err != nil {
					return res, err
				}

				if res.Success {
					callback(e.AppContext)
				} else {
					return res, nil
				}

				return e.Next()
			}),
		}, nil
	case "onAfterItemUse":
		return []event.Unsubscribe{
			ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*event.Result, error) {
				if e.InvItemId != ctx.InvItemID {
					return e.Next()
				}

				res, err := ef.tryToApplyEffect(e.AppContext, ctx.User)
				if err != nil {
					return res, err
				}

				if res.Success {
					callback(e.AppContext)
				} else {
					return res, nil
				}

				return e.Next()
			}),
		}, nil
	default:
		return nil, nil
	}
}

func (ef *GoToJailEffect) tryToApplyEffect(ctx adventuria.AppContext, user adventuria.User) (*event.Result, error) {
	user.SetIsInJail(true)

	_, err := user.MoveToClosestCellType(ctx, "jail")
	if err != nil {
		return &event.Result{
			Success: false,
			Error:   "internal error: can't move to jail cell",
		}, fmt.Errorf("goToJailEffect: %w", err)
	}

	return &event.Result{Success: true}, nil
}

func (ef *GoToJailEffect) Verify(_ adventuria.AppContext, value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

type GoToJailEffectValue struct {
	Event string
}

func (ef *GoToJailEffect) DecodeValue(value string) (*GoToJailEffectValue, error) {
	events := []string{"onAfterItemSave", "onAfterItemUse"}

	if !slices.Contains(events, value) {
		return nil, fmt.Errorf("goToJailEffect: invalid event: %s", value)
	}

	return &GoToJailEffectValue{Event: value}, nil
}

func (ef *GoToJailEffect) GetVariants(_ adventuria.AppContext, _ adventuria.EffectContext) any {
	return nil
}

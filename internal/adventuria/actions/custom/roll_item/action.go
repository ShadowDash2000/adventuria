package roll_item

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_events/custom/item_received"
	"adventuria/pkg/helper"
	"context"
	"errors"
)

type actionsService interface {
	CanDo(ctx context.Context, actor *player.Actor, t model.ActionType) bool
}

type inventories interface {
	AddItemByID(ctx context.Context, actor *player.Actor, itemId string, source item_received.Source) (*effects.InventoryItemWithEffects, error)
	HasEmptySlots(ctx context.Context, player *model.Player) (bool, error)
}

type items interface {
	GetAllRollable(ctx context.Context) ([]*model.Item, error)
}

var _ actions.Action = (*RollItem)(nil)

const Type model.ActionType = "roll_item"

type RollItem struct {
	actions.ActionBase
	actions     actionsService
	inventories inventories
	items       items
}

func NewDef(actionsService actionsService, inventories inventories, items items) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() actions.Action {
			return &RollItem{
				ActionBase:  actions.NewActionBase(Type),
				actions:     actionsService,
				inventories: inventories,
				items:       items,
			}
		},
	)
}

func (r *RollItem) CanDo(ctx context.Context, actor *player.Actor) bool {
	if actor.Player().Progress().ItemWheelsCount() <= 0 {
		return false
	}

	if r.actions.CanDo(ctx, actor, actions.ActionTypeDone) {
		return false
	}

	return true
}

func (r *RollItem) Do(ctx context.Context, actor *player.Actor, _ actions.Request) (any, error) {
	hasEmptySlots, err := r.inventories.HasEmptySlots(ctx, actor.Player())
	if err != nil {
		return nil, err
	}
	if !hasEmptySlots {
		return nil, errs.ErrNotEnoughInventorySlots
	}

	res := model.WheelRollResult{}

	items, err := r.items.GetAllRollable(ctx)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errors.New("no items to roll")
	}

	res.WinnerId = helper.RandomItemFromSlice(items).ID()

	_, err = r.inventories.AddItemByID(ctx, actor, res.WinnerId, item_received.SourceItemWheel)
	if err != nil {
		return nil, err
	}

	err = actor.Player().Progress().ItemWheelsCountChange(-1)
	if err != nil {
		return nil, err
	}

	err = actor.Events().OnAfterItemRoll().Trigger(ctx, &model.OnAfterItemRollEvent{
		ItemId: res.WinnerId,
	})
	if err != nil {
		return nil, err
	}
	err = actor.Events().OnAfterWheelRoll().Trigger(ctx, &model.OnAfterWheelRollEvent{
		ItemId: res.WinnerId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

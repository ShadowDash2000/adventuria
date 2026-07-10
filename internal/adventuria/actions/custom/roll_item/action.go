package roll_item

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/helper"
	"context"
	"errors"
)

type actionsService interface {
	CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool
}

type inventories interface {
	AddItemByID(ctx context.Context, events *model.Events, player *model.Player, itemId string) (*model.InventoryItem, error)
}

type items interface {
	GetAllRollable(ctx context.Context) ([]*model.Item, error)
}

var _ model.Action = (*RollItem)(nil)

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
		func() model.Action {
			return &RollItem{
				ActionBase:  actions.NewActionBase(Type),
				actions:     actionsService,
				inventories: inventories,
				items:       items,
			}
		},
	)
}

func (r *RollItem) CanDo(ctx context.Context, events *model.Events, player *model.Player) bool {
	if player.Progress().ItemWheelsCount() <= 0 {
		return false
	}

	if r.actions.CanDo(ctx, events, player, actions.ActionTypeDone) {
		return false
	}

	return true
}

func (r *RollItem) Do(ctx context.Context, events *model.Events, player *model.Player, _ model.ActionRequest) (any, error) {
	res := model.WheelRollResult{}

	items, err := r.items.GetAllRollable(ctx)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errors.New("no items to roll")
	}

	res.WinnerId = helper.RandomItemFromSlice(items).ID()
	_, err = r.inventories.AddItemByID(ctx, events, player, res.WinnerId)
	if err != nil {
		return nil, err
	}

	err = player.Progress().ItemWheelsCountChange(-1)
	if err != nil {
		return nil, err
	}

	err = events.OnAfterItemRoll().Trigger(ctx, &model.OnAfterItemRollEvent{
		ItemId: res.WinnerId,
	})
	if err != nil {
		return nil, err
	}
	err = events.OnAfterWheelRoll().Trigger(ctx, &model.OnAfterWheelRollEvent{
		ItemId: res.WinnerId,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

package buy

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_events/custom/item_received"
	"context"
	"errors"
	"fmt"
	"slices"
)

type items interface {
	GetByID(ctx context.Context, id string) (*model.Item, error)
	GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error)
}

type inventories interface {
	TryAddItem(ctx context.Context, actor *player.Actor, item *model.Item, source item_received.Source) (*effects.InventoryItemWithEffects, error)
}

var _ actions.Action = (*Buy)(nil)

const Type model.ActionType = "buy"

type Buy struct {
	actions.ActionBase
	items       items
	inventories inventories
}

func NewDef(items items, inventories inventories) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() actions.Action {
			return &Buy{
				ActionBase:  actions.NewActionBase(Type),
				items:       items,
				inventories: inventories,
			}
		},
	)
}

func (b *Buy) CanDo(_ context.Context, actor *player.Actor) bool {
	return actor.Player().LastAction().State().Shop != nil
}

func (b *Buy) Do(ctx context.Context, actor *player.Actor, actionReq actions.Request) (any, error) {
	req, ok := actionReq.(Request)
	if !ok {
		return nil, errors.New("invalid request")
	}
	err := req.validate()
	if err != nil {
		return nil, err
	}

	shopState := actor.Player().LastAction().State().Shop
	if shopState == nil {
		return nil, errs.ErrNoActiveShop
	}

	itemIds := shopState.Ids
	index := slices.Index(itemIds, req.ItemId)
	if index == -1 {
		return nil, fmt.Errorf("item with id = %s not found", req.ItemId)
	}
	itemIds = slices.Delete(itemIds, index, index+1)

	item, err := b.items.GetByID(ctx, req.ItemId)
	if err != nil {
		return nil, err
	}

	buyResult, err := calculatePrice(ctx, actor.Events(), item, shopState, model.EffectRunModeApply)
	if err != nil {
		return nil, err
	}

	if actor.Player().Progress().Balance() < buyResult.Price() {
		return nil, errs.ErrNotEnoughMoney
	}

	_, err = b.inventories.TryAddItem(ctx, actor, item, item_received.SourceShop)
	if err != nil {
		return nil, err
	}

	shopState.Ids = itemIds

	err = actor.Player().Progress().BalanceChange(-buyResult.Price())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

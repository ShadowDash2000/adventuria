package buy

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
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
	TryAddItem(ctx context.Context, events *model.Events, player *model.Player, item *model.Item) (*model.InventoryItem, error)
}

var _ model.Action = (*Buy)(nil)

const Type model.ActionType = "buy"

type Buy struct {
	actions.ActionBase
	items       items
	inventories inventories
}

func NewDef(items items, inventories inventories) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &Buy{
				ActionBase:  actions.NewActionBase(Type),
				items:       items,
				inventories: inventories,
			}
		},
	)
}

func (b *Buy) CanDo(_ context.Context, _ *model.Events, player *model.Player) bool {
	return player.LastAction().State().Shop != nil
}

type Request struct {
	ItemId string `json:"item_id"`
}

func (b *Buy) Do(ctx context.Context, events *model.Events, player *model.Player, actionReq model.ActionRequest) (any, error) {
	req, ok := actionReq.(Request)
	if !ok {
		return nil, errors.New("invalid request")
	}
	if req.ItemId == "" {
		return nil, errors.New("item id is required")
	}

	shopState := player.LastAction().State().Shop
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

	buyResult, err := calculatePrice(ctx, events, item, shopState, model.EffectRunModeApply)
	if err != nil {
		return nil, err
	}

	if player.Progress().Balance() < buyResult.Price() {
		return nil, errs.ErrNotEnoughMoney
	}

	_, err = b.inventories.TryAddItem(ctx, events, player, item)
	if err != nil {
		return nil, err
	}

	shopState.Ids = itemIds

	err = player.Progress().BalanceChange(-buyResult.Price())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

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

type cells interface {
	GetByPlayerWrapped(ctx context.Context, player *model.Player) (model.Cell, error)
}

type items interface {
	GetByID(ctx context.Context, id string) (*model.Item, error)
	GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error)
}

type inventories interface {
	AddItem(ctx context.Context, events *model.Events, player *model.Player, item *model.Item) (*model.InventoryItem, error)
}

var _ model.Action = (*Buy)(nil)

const Type model.ActionType = "buy"

type Buy struct {
	actions.ActionBase
	cells       cells
	items       items
	inventories inventories
}

func NewDef(cells cells, items items, inventories inventories) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &Buy{
				ActionBase:  actions.NewActionBase(Type),
				cells:       cells,
				items:       items,
				inventories: inventories,
			}
		},
	)
}

func (b *Buy) CanDo(ctx context.Context, _ *model.Events, player *model.Player) bool {
	currentCell, err := b.cells.GetByPlayerWrapped(ctx, player)
	if err != nil {
		return false
	}

	if !currentCell.InCategory("shop") {
		return false
	}

	return true
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

	actionState := player.LastAction().State()
	shopState := actionState.Shop
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

	basePrice, err := b.calculatePrice(item.Price(), shopState)
	if err != nil {
		return nil, err
	}

	onBeforeItemBuy, err := b.triggerOnBeforeItemBuy(ctx, events, item, basePrice)
	if err != nil {
		return nil, err
	}

	if player.Progress().Balance() < onBeforeItemBuy.Price {
		return nil, errs.ErrNotEnoughMoney
	}

	_, err = b.inventories.AddItem(ctx, events, player, item)
	if err != nil {
		return nil, err
	}

	actionState.Shop.Ids = itemIds
	player.LastAction().SetState(actionState)

	err = player.Progress().BalanceChange(-onBeforeItemBuy.Price)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

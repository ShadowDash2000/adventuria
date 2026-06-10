package buy

import (
	"adventuria/internal/adventuria_new/actions"
	"adventuria/internal/adventuria_new/errs"
	"adventuria/internal/adventuria_new/model"
	"context"
	"errors"
	"fmt"
	"slices"
)

type cells interface {
	GetCurrentCellByProgress(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error)
}

type items interface {
	GetByID(ctx context.Context, id string) (*model.Item, error)
	GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error)
}

type inventories interface {
	AddItem(ctx context.Context, events *model.Events, playerId string, item *model.Item) (*model.InventoryItem, error)
}

var _ model.Action = (*Buy)(nil)

const Type model.ActionType = "buy"

type Buy struct {
	actions.ActionBase
	cells       cells
	items       items
	inventories inventories
}

func NewActionBuyDef(cells cells, items items, inventories inventories) actions.ActionDef {
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
	currentCell, err := b.cells.GetCurrentCellByProgress(ctx, player.Progress())
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

	ids := player.LastAction().ItemsList()
	index := slices.Index(ids, req.ItemId)
	if index == -1 {
		return nil, fmt.Errorf("item with id = %s not found", req.ItemId)
	}
	ids = slices.Delete(ids, index, index+1)

	item, err := b.items.GetByID(ctx, req.ItemId)
	if err != nil {
		return nil, err
	}

	currentCell, err := b.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return nil, err
	}

	cellShopValue, err := b.decodeValue(currentCell.Data().Value())
	if err != nil {
		return nil, err
	}

	basePrice, err := b.calculatePrice(item.Price(), cellShopValue)
	if err != nil {
		return nil, err
	}

	onBeforeItemBuy, err := b.triggerOnBeforeItemBuy(events, item, basePrice)
	if err != nil {
		return nil, err
	}

	if player.Progress().Balance() < onBeforeItemBuy.Price {
		return nil, errs.ErrNotEnoughMoney
	}

	_, err = b.inventories.AddItem(ctx, events, player.ID(), item)
	if err != nil {
		return nil, err
	}

	player.LastAction().SetItemsList(ids)
	err = player.Progress().BalanceChange(-onBeforeItemBuy.Price)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

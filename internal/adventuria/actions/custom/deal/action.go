package deal

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type inventories interface {
	AddItemByID(ctx context.Context, events *model.Events, player *model.Player, itemId string) (*model.InventoryItem, error)
}

var _ model.Action = (*Deal)(nil)

const Type model.ActionType = "deal"

type Deal struct {
	actions.ActionBase
	inventories inventories
}

func NewDef(inventories inventories) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &Deal{
				ActionBase:  actions.NewActionBase(Type),
				inventories: inventories,
			}
		},
	)
}

func (d *Deal) CanDo(_ context.Context, _ *model.Events, player *model.Player) bool {
	return player.LastAction().State().Dealer != nil
}

func (d *Deal) Do(ctx context.Context, events *model.Events, player *model.Player, _ model.ActionRequest) (any, error) {
	deal := player.LastAction().State().Dealer

	if deal == nil {
		return nil, errs.ErrNoActiveDeals
	}

	switch deal.Type {
	case model.DealTypeCoinsForItem:
		return nil, d.doCoinsForItemDeal(ctx, events, player, deal)
	default:
		return nil, errors.New("unknown deal type")
	}
}

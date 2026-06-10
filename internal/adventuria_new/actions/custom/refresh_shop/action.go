package refresh_shop

import (
	"adventuria/internal/adventuria_new/actions"
	"adventuria/internal/adventuria_new/errs"
	"adventuria/internal/adventuria_new/model"
	"context"
	"errors"
)

type cells interface {
	GetCurrentCellByProgress(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error)
}

var _ model.Action = (*RefreshShop)(nil)

const Type model.ActionType = "refresh_shop"

type RefreshShop struct {
	actions.ActionBase
	cells cells
}

func NewActionRefreshShopDef(cells cells) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &RefreshShop{
				ActionBase: actions.NewActionBase(Type),
				cells:      cells,
			}
		},
	)
}

func (r *RefreshShop) CanDo(ctx context.Context, _ *model.Events, player *model.Player) bool {
	currentCell, err := r.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return false
	}

	if !currentCell.InCategory("shop") {
		return false
	}

	if _, ok := currentCell.(model.Refreshable); !ok {
		return false
	}

	return true
}

func (r *RefreshShop) Do(ctx context.Context, events *model.Events, player *model.Player, _ model.ActionRequest) (any, error) {
	currentCell, err := r.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return nil, err
	}

	cellShopRefreshValue, err := r.decodeValue(currentCell.Data().Value())
	if err != nil {
		return nil, err
	}

	if player.Progress().Balance() < cellShopRefreshValue.RefreshPrice {
		return nil, errs.ErrNotEnoughMoney
	}

	cellRefreshable, ok := currentCell.(model.Refreshable)
	if !ok {
		return nil, errors.New("current cell is not refreshable")
	}

	err = cellRefreshable.RefreshItems(ctx, events, player)
	if err != nil {
		return nil, err
	}

	return nil, player.Progress().BalanceChange(-cellShopRefreshValue.RefreshPrice)
}

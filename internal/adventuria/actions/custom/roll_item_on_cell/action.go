package roll_item_on_cell

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"context"
	"fmt"
)

type cellsService interface {
	GetCurrentCellByProgress(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error)
}

type inventories interface {
	AddItemByID(ctx context.Context, events *model.Events, player *model.Player, itemId string) (*model.InventoryItem, error)
}

type items interface {
	GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error)
}

var _ model.Action = (*RollItemOnCell)(nil)

const Type model.ActionType = "roll_item_on_cell"

type RollItemOnCell struct {
	actions.ActionBase
	cells       cellsService
	inventories inventories
	items       items
}

func NewDef(cells cellsService, inventories inventories, items items) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &RollItemOnCell{
				ActionBase:  actions.NewActionBase(Type),
				cells:       cells,
				inventories: inventories,
				items:       items,
			}
		},
	)
}

func (r *RollItemOnCell) CanDo(ctx context.Context, _ *model.Events, player *model.Player) bool {
	currentCell, err := r.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return false
	}

	if currentCell.Data().Type() != cells.CellTypeRollItem {
		return false
	}

	return !player.Progress().CanMove() && player.LastAction().Type() != Type
}

func (r *RollItemOnCell) Do(ctx context.Context, events *model.Events, player *model.Player, _ model.ActionRequest) (any, error) {
	currentCell, err := r.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return nil, err
	}

	cellRollable, ok := currentCell.(model.Rollable)
	if !ok {
		return nil, fmt.Errorf("current cell is not rollable")
	}

	res, err := cellRollable.Roll(ctx, events, player)
	if err != nil {
		return nil, err
	}

	_, err = r.inventories.AddItemByID(ctx, events, player, res.WinnerId)
	if err != nil {
		return nil, err
	}

	lastAction := player.LastAction()
	lastAction.SetType(Type)

	player.Progress().SetCanMove(true)

	return res, nil
}

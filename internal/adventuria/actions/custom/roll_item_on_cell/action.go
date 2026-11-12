package roll_item_on_cell

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_events/custom/item_received"
	"context"
	"fmt"
)

type cellsService interface {
	GetByPlayer(ctx context.Context, player *model.Player) (*model.CellInfo, error)
	GetByPlayerWrapped(ctx context.Context, player *model.Player) (cells.Cell, error)
}

type inventories interface {
	AddItemByID(ctx context.Context, actor *player.Actor, itemId string, source item_received.Source) (*effects.InventoryItemWithEffects, error)
}

type items interface {
	GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error)
}

var _ actions.Action = (*RollItemOnCell)(nil)

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
		func() actions.Action {
			return &RollItemOnCell{
				ActionBase:  actions.NewActionBase(Type),
				cells:       cells,
				inventories: inventories,
				items:       items,
			}
		},
	)
}

func (r *RollItemOnCell) CanDo(ctx context.Context, actor *player.Actor) bool {
	currentCell, err := r.cells.GetByPlayer(ctx, actor.Player())
	if err != nil {
		return false
	}

	if currentCell.Type() != cells.CellTypeRollItem {
		return false
	}

	return !actor.Player().Progress().CanMove() && actor.Player().LastAction().Status() != model.ActionStatusRollItemOnCell
}

func (r *RollItemOnCell) Do(ctx context.Context, actor *player.Actor, _ actions.Request) (any, error) {
	currentCell, err := r.cells.GetByPlayerWrapped(ctx, actor.Player())
	if err != nil {
		return nil, err
	}

	cellRollable, ok := currentCell.(cells.Rollable)
	if !ok {
		return nil, fmt.Errorf("current cell is not rollable")
	}

	res, err := cellRollable.Roll(ctx, actor.Events(), actor.Player())
	if err != nil {
		return nil, err
	}

	_, err = r.inventories.AddItemByID(ctx, actor, res.WinnerId, item_received.SourceItemWheel)
	if err != nil {
		return nil, err
	}

	lastAction := actor.Player().LastAction()
	lastAction.SetStatus(model.ActionStatusRollItemOnCell)

	actor.Player().Progress().SetCanMove(true)

	return res, nil
}

package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/result"
	"errors"
	"fmt"
)

type RollItemOnCellAction struct {
	adventuria.ActionBase
}

func (a *RollItemOnCellAction) CanDo(ctx adventuria.ActionContext) bool {
	currentCell, ok := ctx.User.CurrentCell()
	if !ok {
		return false
	}

	if currentCell.Type() != "rollItem" {
		return false
	}

	return !ctx.User.LastAction().CanMove() && ctx.User.LastAction().Type() != ActionTypeRollItemOnCell
}

func (a *RollItemOnCellAction) Do(ctx adventuria.ActionContext, req adventuria.ActionRequest) (*result.Result, error) {
	currentCell, ok := ctx.User.CurrentCell()
	if !ok {
		return result.Err("internal error: current cell not found"),
			errors.New("roll_item_on_cell.do(): current cell not found")
	}

	cellWheel := currentCell.(adventuria.CellWheel)

	res, err := cellWheel.Roll(ctx.AppContext, ctx.User, adventuria.RollWheelRequest(req))
	if err != nil {
		return result.Err("internal error: failed to roll an item"),
			fmt.Errorf("roll_item_on_cell.do(): %w", err)
	}

	_, err = ctx.User.Inventory().MustAddItemById(ctx.AppContext, res.WinnerId)
	if err != nil {
		return result.Err("internal error: failed to add item to the inventory"),
			fmt.Errorf("roll_item_on_cell.do(): can't add item to inventory: %w", err)
	}

	action := ctx.User.LastAction()
	action.SetType(ActionTypeRollItemOnCell)
	action.SetCanMove(true)

	return result.Ok().WithData(res), nil
}

func (a *RollItemOnCellAction) GetVariants(_ adventuria.ActionContext) any {
	return nil
}

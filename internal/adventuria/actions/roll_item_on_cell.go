package actions

import (
	"adventuria/internal/adventuria"
	"errors"
	"fmt"
)

type RollItemOnCellAction struct {
	adventuria.ActionBase
}

func (a *RollItemOnCellAction) CanDo(user adventuria.User) bool {
	currentCell, ok := user.CurrentCell()
	if !ok {
		return false
	}

	if currentCell.Type() != "rollItem" {
		return false
	}

	return !user.LastAction().CanMove() && user.LastAction().Type() != ActionTypeRollItemOnCell
}

func (a *RollItemOnCellAction) Do(user adventuria.User, req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	currentCell, ok := user.CurrentCell()
	if !ok {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: current cell not found",
		}, errors.New("roll_item_on_cell.do(): current cell not found")
	}

	cellWheel := currentCell.(adventuria.CellWheel)

	res, err := cellWheel.Roll(user, adventuria.RollWheelRequest(req))
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error",
		}, fmt.Errorf("roll_item_on_cell.do(): %w", err)
	}

	_, err = user.Inventory().MustAddItemById(res.WinnerId)
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error",
		}, fmt.Errorf("roll_item_on_cell.do(): can't add item to inventory: %w", err)
	}

	action := user.LastAction()
	action.SetType(ActionTypeRollItemOnCell)
	action.SetCanMove(true)

	return &adventuria.ActionResult{
		Success: true,
		Data:    res,
	}, nil
}

package actions

import (
	"adventuria/internal/adventuria"
	"errors"
	"fmt"
)

type RollWheelAction struct {
	adventuria.ActionBase
}

func (a *RollWheelAction) CanDo(user adventuria.User) bool {
	currentCell, ok := user.CurrentCell()
	if !ok {
		return false
	}

	if _, ok = currentCell.(adventuria.CellWheel); !ok {
		return false
	}

	return !user.LastAction().CanMove() && user.LastAction().Type() != ActionTypeRollWheel
}

func (a *RollWheelAction) Do(user adventuria.User, req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	currentCell, ok := user.CurrentCell()
	if !ok {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: current cell not found",
		}, errors.New("roll_wheel.do(): current cell not found")
	}

	onBeforeWheelRollEvent := &adventuria.OnBeforeWheelRollEvent{
		CurrentCell: currentCell.(adventuria.CellWheel),
	}
	err := user.OnBeforeWheelRoll().Trigger(onBeforeWheelRollEvent)
	if err != nil {
		adventuria.PocketBase.Logger().Error(
			"roll_wheel.do(): failed to trigger onBeforeWheelRoll event",
			"error",
			err,
		)
	}

	res, err := onBeforeWheelRollEvent.CurrentCell.Roll(user, adventuria.RollWheelRequest(req))
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error",
		}, fmt.Errorf("roll_wheel.do(): %w", err)
	}

	action := user.LastAction()
	action.SetType(ActionTypeRollWheel)
	action.SetGame(res.WinnerId)

	onAfterWheelRollEvent := &adventuria.OnAfterWheelRollEvent{
		ItemId: res.WinnerId,
	}
	err = user.OnAfterWheelRoll().Trigger(onAfterWheelRollEvent)
	if err != nil {
		adventuria.PocketBase.Logger().Error(
			"roll_wheel.do(): failed to trigger onAfterWheelRoll event",
			"error",
			err,
		)
	}

	return &adventuria.ActionResult{
		Success: true,
		Data:    res,
	}, nil
}

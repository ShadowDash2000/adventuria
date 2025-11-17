package actions

import (
	"adventuria/internal/adventuria"
	"errors"
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

	return !user.LastAction().CanMove()
}

func (a *RollWheelAction) Do(user adventuria.User, req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	currentCell, ok := user.CurrentCell()
	if !ok {
		return nil, errors.New("current cell not found")
	}

	onBeforeWheelRollEvent := &adventuria.OnBeforeWheelRollEvent{
		CurrentCell: currentCell.(adventuria.CellWheel),
	}
	err := user.OnBeforeWheelRoll().Trigger(onBeforeWheelRollEvent)
	if err != nil {
		return nil, err
	}

	res, err := onBeforeWheelRollEvent.CurrentCell.Roll(user, adventuria.RollWheelRequest(req))
	if err != nil {
		return nil, err
	}

	action := user.LastAction()
	action.SetType(ActionTypeRollWheel)
	action.SetGame(res.WinnerId)

	onAfterWheelRollEvent := &adventuria.OnAfterWheelRollEvent{
		ItemId: res.WinnerId,
	}
	err = user.OnAfterWheelRoll().Trigger(onAfterWheelRollEvent)
	if err != nil {
		return nil, err
	}

	return &adventuria.ActionResult{
		Success: true,
		Data:    res,
	}, nil
}

package actions

import (
	"adventuria/internal/adventuria"
	"errors"
)

type RollWheelAction struct {
	adventuria.ActionBase
}

func (a *RollWheelAction) CanDo() bool {
	return a.User().GetNextStepType() == adventuria.ActionTypeRollWheel
}

func (a *RollWheelAction) Do(_ adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	currentCell, ok := a.User().CurrentCell()
	if !ok {
		return nil, errors.New("current cell not found")
	}

	onBeforeWheelRollEvent := &adventuria.OnBeforeWheelRollEvent{
		CurrentCell: currentCell.(adventuria.CellWheel),
	}
	err := a.User().OnBeforeWheelRoll().Trigger(onBeforeWheelRollEvent)
	if err != nil {
		return nil, err
	}

	res, err := onBeforeWheelRollEvent.CurrentCell.Roll(a.User())
	if err != nil {
		return nil, err
	}

	action := a.User().LastAction()
	action.SetType(adventuria.ActionTypeRollWheel)
	action.SetValue(res.WinnerId)
	action.SetCollectionRef(res.Collection.Id)

	onAfterWheelRollEvent := &adventuria.OnAfterWheelRollEvent{
		ItemId: res.WinnerId,
	}
	err = a.User().OnAfterWheelRoll().Trigger(onAfterWheelRollEvent)
	if err != nil {
		return nil, err
	}

	return &adventuria.ActionResult{
		Success: true,
		Data:    res,
	}, nil
}

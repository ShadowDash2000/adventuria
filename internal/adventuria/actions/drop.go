package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"errors"
)

type DropAction struct {
	adventuria.ActionBase
}

func (a *DropAction) CanDo() bool {
	return a.User().GetNextStepType() == adventuria.ActionTypeChooseResult
}

func (a *DropAction) Do(req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	currentCell, ok := a.User().CurrentCell()
	if !ok {
		return nil, errors.New("current cell not found")
	}

	if currentCell.CantDrop() || !a.User().CanDrop() {
		return nil, errors.New("can't drop on this cell")
	}

	onBeforeDropEvent := &adventuria.OnBeforeDropEvent{
		IsSafeDrop: false,
	}
	err := a.User().OnBeforeDrop().Trigger(onBeforeDropEvent)
	if err != nil {
		return nil, err
	}

	action := a.User().LastAction()
	action.SetType(adventuria.ActionTypeDrop)
	action.SetComment(req.Comment)

	if !onBeforeDropEvent.IsSafeDrop && !currentCell.IsSafeDrop() {
		a.User().SetPoints(a.User().Points() + adventuria.GameSettings.PointsForDrop())
		a.User().SetDropsInARow(a.User().DropsInARow() + 1)

		if !a.User().IsSafeDrop() {
			if err = a.moveToJail(); err != nil {
				return nil, err
			}
		}
	}

	err = a.User().OnAfterDrop().Trigger(&adventuria.OnAfterDropEvent{})
	if err != nil {
		return nil, err
	}

	return &adventuria.ActionResult{
		Success: true,
	}, nil
}

func (a *DropAction) moveToJail() error {
	err := a.User().MoveToCellType(cells.CellTypeJail)
	if err != nil {
		return err
	}

	err = a.User().OnAfterGoToJail().Trigger(&adventuria.OnAfterGoToJailEvent{})
	if err != nil {
		return err
	}

	return nil
}

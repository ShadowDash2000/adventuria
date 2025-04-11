package adventuria

import (
	"adventuria/pkg/helper"
	"errors"
)

type CellJail struct {
	CellWheelBase
}

func NewCellJail() CellCreator {
	return func() Cell {
		return &CellJail{}
	}
}

func (c *CellJail) NextStep(user *User) string {
	nextStepType := ""

	if user.IsInJail() {
		switch user.LastAction.Type() {
		case ActionTypeRoll,
			ActionTypeReroll,
			ActionTypeDrop:
			nextStepType = ActionTypeRollCell
		case ActionTypeRollCell:
			nextStepType = ActionTypeChooseGame
		case ActionTypeChooseGame:
			nextStepType = ActionTypeDone
		case ActionTypeDone:
			nextStepType = ActionTypeRoll
		default:
			nextStepType = ActionTypeRollCell
		}
	} else {
		nextStepType = ActionTypeRoll
	}

	return nextStepType
}

func (c *CellJail) Roll(_ *User) (*WheelRollResult, error) {
	cellsCol, err := c.gc.Cols.Get(TableCells)
	if err != nil {
		return nil, err
	}

	res := &WheelRollResult{
		Collection: cellsCol,
	}

	cells := c.gc.Cells.GetAllByType(CellTypeGame)

	if len(cells) == 0 {
		return nil, errors.New("game cells not found")
	}

	for _, cell := range cells {
		res.FillerItems = append(res.FillerItems, &WheelItem{
			Name: cell.Name(),
			Icon: cell.Icon(),
		})
	}

	res.WinnerId = helper.RandomItemFromSlice(cells).ID()

	return res, nil
}
